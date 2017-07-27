package dns

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"time"
)

type queryHandler interface {
	Handle(*mdns.Msg, *config.Zone, mdns.Question) error
}

type QueryHandlers map[string]queryHandler

type DNS struct {
	Config   *config.Config
	Server   *mdns.Server
	handlers QueryHandlers
}

var now = func() int64 {
	return time.Now().Unix()
}

// NewDNS created new instance of DNS server
func NewDNS(config *config.Config) *DNS {
	h := QueryHandlers{
		"A":   &aHandler{},
		"MX":  &mxHandler{},
		"TXT": &txtHandler{},
		"NS":  NewNSHandler(config),
	}

	if len(config.DNS.Servers) > 0 {
		h["SOA"] = NewSOAHandler(config.DNS.Servers[0].Name)
	}

	dns := &DNS{
		Config:   config,
		handlers: h,
	}
	return dns
}

// HandleRequest process incoming requests
func (dns *DNS) HandleRequest(w mdns.ResponseWriter, r *mdns.Msg) {
	msg := &mdns.Msg{}
	msg.SetReply(r)
	msg.Compress = false
	msg.Authoritative = true

	switch r.Opcode {
	case mdns.OpcodeQuery:
		dns.ParseQuery(msg)
	}

	w.WriteMsg(msg)
}

// Serve starts DNS server
func (dns *DNS) Serve(net string) (err error) {
	log.WithFields(log.Fields{
		"service": "dns",
		"net":     net,
	}).Info("Serve over ", net)

	dns.Server = &mdns.Server{
		Addr: fmt.Sprintf("%s:%d", dns.Config.DNS.IP, dns.Config.DNS.Port),
		Net:  net,
	}

	mdns.HandleFunc(".", dns.HandleRequest)

	err = dns.Server.ListenAndServe()
	if err != nil {
		log.WithFields(log.Fields{
			"service": "dns",
			"addr":    dns.Config.DNS.IP,
			"port":    dns.Config.DNS.Port,
			"net":     net,
			"error":   err,
		}).Error("Can't start server")
	}
	defer dns.Server.Shutdown()

	log.WithFields(log.Fields{
		"service": "dns",
	}).Info("[dns] Shutdown...")

	return
}
