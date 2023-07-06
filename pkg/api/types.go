package api

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HandlerIFace interface {
	GetRoutes() []Route
	GetService() Service
	GetPrefix() string
}

type Handler struct {
	HandlerIFace
	Routes  []Route
	Prefix  string
	Service Service
}

type Route struct {
	Name        string           `json:"name"`
	Method      string           `json:"method"`
	Pattern     string           `json:"pattern"`
	Protected   bool             `json:"protected"`
	HandlerFunc http.HandlerFunc `json:"-"`
}

type Service interface {
	GetName() (name string)
	//GetStore() (store DataStoreIFace)
}

type DataService struct {
	Service
	name string
}

func NewDataService() DataService {
	return DataService{
		name: "Data",
	}
}

func (ds DataService) GetName() (name string) {
	return ds.name
}

func (ds DataService) Write() (err error) {
	log.Info("In Service Write")
	return nil
}
