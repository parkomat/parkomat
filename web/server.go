package web

import (
	"fmt"
	"github.com/golang/glog"
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
	glog.Info("[web] Serve...")

	var log *os.File = os.Stdout
	if server.Config.Web.AccessLog != "" {
		log, err = os.OpenFile(server.Config.Web.AccessLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			glog.Error("[web] Can't create log file.", err)
			return
		}
	}

	hl := NewHttpLogHandler(server.mux, log)

	err = http.ListenAndServe(
		fmt.Sprintf("%s:%d",
			server.Config.Web.IP,
			server.Config.Web.Port),
		hl)

	if err != nil {
		glog.Error("[web] Listen error: ", err)
		return
	}

	glog.Info("[web] Shutdown...")
	return
}
