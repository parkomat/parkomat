package web

import (
	"fmt"
	"github.com/golang/glog"
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

		glog.Info("[web] Request: ", r.Host, " ", r.URL.Path)

		// TODO: optimize this monster! ;)
		fp := path.Join(server.Config.Web.Path, r.Host, "public_html", r.URL.Path)
		_, err := os.Stat(fp)
		if err != nil {
			fp = path.Join(server.Config.Web.Path, "default", "public_html", r.URL.Path)
			_, err = os.Stat(fp)
			if err != nil {
				if os.IsNotExist(err) {
					glog.Error("[web] Not found: ", fp)
					http.NotFound(w, r)
					return
				}
			}
		}

		http.ServeFile(w, r, fp)
	}
}
