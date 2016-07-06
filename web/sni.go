package web

/*
	Inspired by:

	https://groups.google.com/forum/#!msg/golang-nuts/rUm2iYTdrU4/PaEBya4dzvoJ
*/

import (
	"crypto/tls"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
)

func (server *Server) ListenAndServeTLSSNI() (err error) {
	// TODO: fix it!
	var lf *os.File = os.Stdout
	if server.Config.Web.AccessLog != "" {
		lf, err = os.OpenFile(server.Config.Web.AccessLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.WithFields(log.Fields{
				"service": "web",
				"path":    server.Config.Web.AccessLog,
				"error":   err,
			}).Error("Can't create log file")
		}
	}
	hl := NewHttpLogHandler(server.mux, lf)

	hs := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", server.Config.Web.IP, server.Config.Web.SSLPort),
		Handler: hl,
	}

	config := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{},
	}

	log.WithFields(log.Fields{
		"service": "web",
		"path":    server.Config.Web.Path,
	}).Info("Loading certificates")

	// Let's walk through our www dir hoping to find some certificates
	files, _ := ioutil.ReadDir(server.Config.Web.Path)
	for _, f := range files {
		if f.Name() != "default" {
			key := path.Join(server.Config.Web.Path, f.Name(), fmt.Sprintf("%s.key", f.Name()))
			crt := path.Join(server.Config.Web.Path, f.Name(), fmt.Sprintf("%s.crt", f.Name()))

			_, err1 := os.Stat(key)
			if err1 != nil {
				log.WithFields(log.Fields{
					"service": "web",
					"name":    f.Name(),
				}).Warning("No .key file")
				continue
			}
			_, err1 = os.Stat(crt)
			if err1 != nil {
				log.WithFields(log.Fields{
					"service": "web",
					"name":    f.Name(),
				}).Warning("No .crt file")
				continue
			}

			log.WithFields(log.Fields{
				"service": "web",
				"name":    f.Name(),
			}).Info("Adding SSL certificate")

			cert, err1 := tls.LoadX509KeyPair(crt, key)
			if err1 != nil {
				log.WithFields(log.Fields{
					"service": "web",
					"name":    f.Name(),
					"error":   err1,
				}).Error("Invalid SSL certificate")
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
		log.WithFields(log.Fields{
			"service": "web",
			"addr":    hs.Addr,
			"error":   err,
		}).Error("Can't listen")
		return
	}

	tlsListener := tls.NewListener(conn, config)
	return hs.Serve(tlsListener)
}
