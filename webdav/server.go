package webdav

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/parkomat/parkomat/config"
	"golang.org/x/net/webdav"
	"net/http"
	"path/filepath"
)

type WebDav struct {
	Config *config.Config

	fs webdav.FileSystem

	HandlerFunc http.HandlerFunc
}

func (wd *WebDav) mount(path string) error {
	if s, err := filepath.Abs(path); err == nil {
		path = s
	}
	wd.fs = webdav.Dir(path)
	return nil
}

func NewWebDav(config *config.Config) *WebDav {
	return &WebDav{
		Config: config,
	}
}

func (wd *WebDav) Init() error {
	if !(wd.Config != nil && wd.Config.WebDav.Enabled != false) {
		return fmt.Errorf("WebDav not configured")
	}

	log.WithFields(log.Fields{
		"service": "webdav",
	}).Info("Init")

	wd.mount(wd.Config.Web.Path)

	h := &webdav.Handler{
		Prefix:     wd.Config.WebDav.Mount,
		FileSystem: wd.fs,
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			log.WithFields(log.Fields{
				"service": "webdav",
				"method":  r.Method,
				"path":    r.URL.Path,
				"error":   err,
			}).Info("Request")
		},
	}

	wd.HandlerFunc = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !(ok == true && u == wd.Config.WebDav.Username && p == wd.Config.WebDav.Password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="davfs"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})

	return nil
}
