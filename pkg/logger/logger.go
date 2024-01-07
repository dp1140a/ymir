package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Logger struct {
	Log    *log.Logger
	Config *LoggerConfig
}

var (
	logger   *Logger
	logLevel = map[string]log.Level{
		"DEBUG": log.DebugLevel,
		"INFO":  log.InfoLevel,
		"WARN":  log.WarnLevel,
		"ERROR": log.ErrorLevel,
		"FATAL": log.FatalLevel,
	}
)

func init() {
	logger = NewLogger()
}

func NewLogger() *Logger {
	var l = new(Logger)
	l.Log = log.New()
	return l
}

func InitLogger() (err error) {
	var sol io.Writer
	var f io.Writer
	logger.Config = NewLoggerConfig()

	//Set Log Level
	log.SetLevel(logLevel[logger.Config.LogLevel])

	// Set Log Formatter
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:   time.RFC3339,
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			repopath := fmt.Sprintf("%s/src/github.com/bob", os.Getenv("GOPATH"))
			filename := strings.Replace(f.File, repopath, "", -1)
			r, _ := regexp.Compile(`[^\/]+\/[^\/]+$`)
			daFunc := strings.Split(f.Function, ".")
			return "", fmt.Sprintf("%s:%d[%s()]", r.Find([]byte(filename)), f.Line, daFunc[len(daFunc)-1])
		},
		PrettyPrint: false,
	})
	log.SetReportCaller(true)

	if logger.Config.StdOut != true {
		sol = io.Discard
	} else {
		sol = os.Stdout
		log.Info("console log enabled")
	}

	if logger.Config.FileOut == true {
		dir, _ := filepath.Split(logger.Config.LogFile)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				log.Fatal("Log Error: ", err, ": ", logger.Config.LogFile)
				return err
			}
		}
		log.Info("http file log enabled")
		log.Info("http log file name: ", logger.Config.LogFile)
		f, err = os.OpenFile(logger.Config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal("Log Error: ", err, " ", logger.Config.LogFile)
			return err
		}
	} else {
		f = io.Discard
	}

	mw := io.MultiWriter(sol, f)
	log.SetOutput(mw)
	log.Info("Logging started")

	return nil
}
