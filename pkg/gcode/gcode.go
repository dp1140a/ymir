package gcode

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type GCode struct {
	MetaData GCodeMetaData
	FilePath string
}

type GCodeMetaData struct {
	GCodeType      string `json:"gCodeType,omitempty"`
	CreatedBy      string `json:"createdBy,omitempty"`
	CreatedDate    string `json:"createDate,omitempty"`
	TotalTime      string `json:"totalTime,omitempty"`
	LayerHeight    string `json:"layerHeight,omitempty"`
	NozzleDiameter string `json:"nozzleDiameter,omitempty"`
	Material       string `json:"material,omitempty"`
	FilamentUsedG  string `json:"filamentUsedG,omitempty"`
	FilamentUsedM  string `json:"filamentUsedM,omitempty"`
	PrinterType    string `json:"printerType,omitempty"`
	Thumbnail      string `json:"thumbnail,omitempty"`
}

func NewGCode(filePath string) *GCode {
	return &GCode{
		MetaData: GCodeMetaData{},
		FilePath: filePath,
	}
}

func (gc *GCode) ParseGCode(debug bool) error {
	log.Infof("Parsing Gcode file: %v", gc.FilePath)
	// first open the file
	file, err := os.Open(gc.FilePath)
	if err != nil {
		log.Errorf("could not open the file: %v", err)
		return err
	}
	// don't forget to close the file.
	defer file.Close()
	// finally, we can have our scanner
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	gCodeType := GetGCodeType(line)

	if strings.HasPrefix(line, ";") {
		if gCodeType == "MARLIN" {
			gc.ParseMarlin(scanner)
		} else if gCodeType == "PRUSA" {
			gc.ParsePrusa(scanner)
		} else {
			//log.Errorf("Unknown GCode Type")
			return errors.New("unknown GCode Type")
		}
	}

	if debug {
		bytes, _ := json.MarshalIndent(gc.MetaData, "", "\t")
		fmt.Println(string(bytes))
	}

	return nil
}

func (gc *GCode) ParsePrusa(scanner *bufio.Scanner) error {
	lineNumber := 1
	gc.MetaData.GCodeType = "PRUSA"
	line := scanner.Text()[1:]
	if strings.Contains(line, "generated by") {
		str := strings.Split(line, "on")
		gc.MetaData.CreatedBy = str[0]
		ts := strings.Split(str[1], "at")
		gc.MetaData.CreatedDate = fmt.Sprintf("%v %v", strings.TrimSpace(ts[0]), strings.TrimSpace(ts[1]))
	}
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if strings.HasPrefix(line, ";") { //Its a Comment Line
			line = strings.TrimSpace(line[1:])
			if strings.HasPrefix(line, "thumbnail begin") {
				thumbNail, endLine := ExtractThumbnail(scanner, lineNumber)
				lineNumber = endLine
				gc.MetaData.Thumbnail = fmt.Sprintf("data:image/png;base64,%s", thumbNail)
			} else {
				kv := strings.Split(line, "=")
				switch strings.TrimSpace(kv[0]) {
				case "estimated printing time (normal mode)":
					gc.MetaData.TotalTime = strings.TrimSpace(kv[1])
				case "layer_height":
					gc.MetaData.LayerHeight = strings.TrimSpace(kv[1])
				case "total filament used [g]":
					gc.MetaData.FilamentUsedG = strings.TrimSpace(kv[1])
				case "filament_type":
					gc.MetaData.Material = strings.TrimSpace(kv[1])
				case "nozzle_diameter":
					gc.MetaData.NozzleDiameter = strings.TrimSpace(kv[1])
				case "filament used [mm]":
					gc.MetaData.FilamentUsedM = strings.TrimSpace(kv[1])
				case "printer_model":
					gc.MetaData.PrinterType = strings.TrimSpace(kv[1])
				default:
					continue
				}
			}

		}
		// else contnue
		if err := scanner.Err(); err != nil {
			return errors.New(fmt.Sprintf("error scanning %v line %v: %v", gc.FilePath, lineNumber, err))
		}
	}
	return nil
}

func B64ToImg(b64Str string) (thumbnail image.Image, err error) {
	unbased, err := base64.StdEncoding.DecodeString(b64Str)
	if err != nil {
		return nil, errors.New("\"Cannot decode b64\"")
	}
	r := bytes.NewReader(unbased)
	thumbnail, err = png.Decode(r)
	if err != nil {
		return nil, errors.New("not a png")
	}
	return
}

func ExtractThumbnail(scanner *bufio.Scanner, startLine int) (b64Str string, endline int) {
	var sb strings.Builder
	for scanner.Scan() {
		startLine++
		line := strings.TrimSpace(scanner.Text()[1:])
		if strings.HasPrefix(line, "thumbnail end") {
			break
		} else {
			sb.WriteString(line)
		}
	}

	return sb.String(), startLine
}

func (gc *GCode) ParseMarlin(scanner *bufio.Scanner) error {
	lineNumber := 1
	gc.MetaData.GCodeType = "MARLIN"
	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++
		if strings.HasPrefix(line, ";") { //Its a Comment Line
			line = line[1:]
			if strings.Contains(line, "Generated") {
				gc.MetaData.CreatedBy = line

			} else {
				kv := strings.Split(line, ":")
				switch kv[0] {
				case "TIME":
					tInt, err := strconv.Atoi(kv[1])
					if err != nil {
						log.Errorf(err.Error())
					}
					gc.MetaData.TotalTime = (time.Duration(tInt) * time.Second).String()
				case "Filament used":
					gc.MetaData.FilamentUsedM = kv[1]
				case "Layer height":
					gc.MetaData.LayerHeight = kv[1]
				default:
					continue
				}
			}

		}
		// else contnue
		if err := scanner.Err(); err != nil {
			return errors.New(fmt.Sprintf("error scanning %v line %v: %v", gc.FilePath, lineNumber, err))
		}
	}

	return nil
}

func GetGCodeType(line string) string {
	if strings.Contains(line, "PrusaSlicer") {
		return "PRUSA"
	} else if strings.Contains(line, "Marlin") {
		return "MARLIN"
	} else {
		return "UNK"
	}
}
