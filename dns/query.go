package dns

import (
	log "github.com/Sirupsen/logrus"
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
)

// ParseQuery parses incoming query and creates appropriate response message
func (dns *DNS) ParseQuery(msg *mdns.Msg) {
	for _, q := range msg.Question {
		t := mdns.TypeToString[q.Qtype]

		log.WithFields(log.Fields{
			"service": "dns",
			"query":   t,
			"name":    q.Name,
		}).Info("Query")

		d := dns.Config.GetDomain(q.Name)
		if d == nil {
			log.WithFields(log.Fields{
				"service": "dns",
				"name":    q.Name,
			}).Error("Domain not configured.")
			return
		}

		// Check whether to use global zone or not
		var z *config.Zone = nil
		if d.Zone != nil {
			z = d.Zone
		} else {
			z = &dns.Config.Zone
		}

		if handler, exists := dns.handlers[t]; exists {
			err := handler.Handle(msg, z, q)
			if err != nil {
				log.WithFields(log.Fields{
					"service": "dns",
					"query":   t,
					"name":    q.Name,
					"error":   err,
				}).Warning("Problem with query")
			}
		} else {
			log.WithFields(log.Fields{
				"service": "dns",
				"query":   t,
				"name":    q.Name,
			}).Warning("Handler doesn't exist")
		}
	}
	return
}
