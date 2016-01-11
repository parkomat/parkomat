package web

/*
	Inspired by:

	https://groups.google.com/forum/#!msg/golang-nuts/rUm2iYTdrU4/PaEBya4dzvoJ
*/

import (
	"crypto/tls"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
)

func (server *Server) ListenAndServeTLSSNI() (err error) {
	// TODO: fix it!
	var log *os.File = os.Stdout
	if server.Config.Web.AccessLog != "" {
		log, err = os.OpenFile(server.Config.Web.AccessLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			glog.Error("[web] Can't create log file.", err)
		}
	}
	hl := NewHttpLogHandler(server.mux, log)

	hs := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", server.Config.Web.IP, server.Config.Web.SSLPort),
		Handler: hl,
	}

	config := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{},
	}

	glog.Info("[web] Looking for certificates in ", server.Config.Web.Path)

	// Let's walk through our www dir hoping to find some certificates
	files, _ := ioutil.ReadDir(server.Config.Web.Path)
	for _, f := range files {
		glog.Info("[web] ", f.Name())
		if f.Name() != "default" {
			key := path.Join(server.Config.Web.Path, f.Name(), fmt.Sprintf("%s.key", f.Name()))
			crt := path.Join(server.Config.Web.Path, f.Name(), fmt.Sprintf("%s.crt", f.Name()))

			_, err1 := os.Stat(key)
			if err1 != nil {
				glog.Info("[web] No key file for ", f.Name(), " domain.")
				continue
			}
			_, err1 = os.Stat(crt)
			if err1 != nil {
				glog.Info("[web] No crt file for ", f.Name(), " domain.")
				continue
			}

			glog.Info("[web] Adding SSL cert for ", f.Name(), " domain.")

			cert, err1 := tls.LoadX509KeyPair(crt, key)
			if err1 != nil {
				glog.Error("[web] Invalid SSL cert for ", f.Name(), " domain.")
				continue
			}
			d := server.Config.GetDomain(f.Name())
			if d != nil {
				d.HasSSL = true
				d.SSLCertificate = crt
				d.SSLCertificateKey = key
			}

			config.Certificates = append(config.Certificates, cert)
		}
	}

	config.BuildNameToCertificate()

	conn, err := net.Listen("tcp", hs.Addr)
	if err != nil {
		glog.Error("[web] Can't listen: ", err)
		return
	}

	tlsListener := tls.NewListener(conn, config)
	return hs.Serve(tlsListener)
}
