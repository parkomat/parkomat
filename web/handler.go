package web

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"path"
)

func (server *Server) handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d := server.Config.GetDomain(r.Host)
		if d == nil {
			http.NotFound(w, r)
			return
		}

		if d.HasSSL == true && r.TLS == nil {
			http.Redirect(w, r, fmt.Sprintf("https://%s/%s", r.Host, r.URL.Path), 301)
			return
		}

		log.WithFields(log.Fields{
			"service": "web",
			"host":    r.Host,
			"path":    r.URL.Path,
		}).Info("Request")

		// TODO: optimize this monster! ;)
		fp := path.Join(server.Config.Web.Path, r.Host, "public_html", r.URL.Path)
		_, err := os.Stat(fp)
		if err != nil {
			fp = path.Join(server.Config.Web.Path, "default", "public_html", r.URL.Path)
			_, err = os.Stat(fp)
			if err != nil {
				if os.IsNotExist(err) {
					log.WithFields(log.Fields{
						"service": "web",
						"path":    fp,
					}).Warning("Not found")
					http.NotFound(w, r)
					return
				}
			}
		}

		http.ServeFile(w, r, fp)
	}
}
