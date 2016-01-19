package dns

import (
	"github.com/golang/glog"
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"strconv"
)

type DNS struct {
	Config *config.Config

	Server *mdns.Server
}

func NewDNS(config *config.Config) *DNS {
	return &DNS{
		Config: config,
	}
}

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

func (dns *DNS) Serve(net string) (err error) {
	glog.Info("[dns] Serve over ", net, "...")

	dns.Server = &mdns.Server{
		Addr: dns.Config.DNS.IP + ":" + strconv.Itoa(dns.Config.DNS.Port),
		Net:  net,
	}

	mdns.HandleFunc(".", dns.HandleRequest)

	err = dns.Server.ListenAndServe()
	if err != nil {
		glog.Error("[dns] Can't start server: ", err)
	}
	defer dns.Server.Shutdown()

	glog.Info("[dns] Shutdown...")

	return
}
