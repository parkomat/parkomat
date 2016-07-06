package web

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/parkomat/parkomat/config"
	"net/http"
	"os"
)

type Server struct {
	Config *config.Config

	mux *http.ServeMux
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Config: config,
	}
}

func (server *Server) Init() (err error) {
	server.mux = http.DefaultServeMux
	server.mux.HandleFunc("/", server.handler())
	return
}

func (server *Server) AddHandlerFunc(path string, handlerFunc http.HandlerFunc) {
	server.mux.HandleFunc(path, handlerFunc)
}

func (server *Server) Serve() (err error) {
	log.WithFields(log.Fields{
		"service": "web",
		"ip":      server.Config.Web.IP,
		"port":    server.Config.Web.Port,
	}).Info("Serve")

	var lf *os.File = os.Stdout
	if server.Config.Web.AccessLog != "" {
		lf, err = os.OpenFile(server.Config.Web.AccessLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.WithFields(log.Fields{
				"service": "web",
				"path":    server.Config.Web.AccessLog,
				"error":   err,
			}).Error("Can't create log file")
			return
		}
	}

	hl := NewHttpLogHandler(server.mux, lf)

	err = http.ListenAndServe(
		fmt.Sprintf("%s:%d",
			server.Config.Web.IP,
			server.Config.Web.Port),
		hl)

	if err != nil {
		log.WithFields(log.Fields{
			"service": "web",
			"ip":      server.Config.Web.IP,
			"port":    server.Config.Web.Port,
		}).Error("Can't listen")
		return
	}

	log.WithFields(log.Fields{
		"service": "web",
	}).Info("Shutdown")
	return
}
