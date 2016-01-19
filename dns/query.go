package dns

import (
	"fmt"
	"github.com/golang/glog"
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"strings"
	"time"
)

// TODO: split each case into separate functions for cleanliness
func (dns *DNS) ParseQuery(msg *mdns.Msg) {
	for _, q := range msg.Question {
		t := mdns.TypeToString[q.Qtype]

		glog.Info("[dns] Query ", t, ": ", q.Name)

		d := dns.Config.GetDomain(q.Name)
		if d == nil {
			glog.Error("[dns] Domain ", q.Name, " not configured.")
			return
		}

		// Check whether to use global zone or not
		var z *config.Zone = nil
		if d.Zone != nil {
			z = d.Zone
		} else {
			z = &dns.Config.Zone
		}

		switch t {
		case "A":
			s := strings.Join(
				[]string{
					q.Name,
					"3600",
					"IN",
					"A",
					z.A,
				}, " ")

			rr, err := mdns.NewRR(s)
			if err == nil {
				msg.Answer = append(msg.Answer, rr)
			}

		case "MX":
			for _, s := range strings.Split(z.MX, "\n") {
				s = strings.Trim(s, " ")
				if s != "" {
					mx := strings.Split(s, " ")

					s = strings.Join([]string{
						q.Name,
						"3600",
						"IN",
						"MX",
						mx[0],
						mx[1],
					}, " ")

					rr, err := mdns.NewRR(s)
					if err == nil {
						msg.Answer = append(msg.Answer, rr)
					}
				}
			}

		case "SOA":
			s := strings.Join(
				[]string{
					q.Name,
					"3600",
					"IN",
					"SOA",
					dns.Config.DNS.Servers[0].Name,
					fmt.Sprintf("admin.%s", q.Name),
					fmt.Sprintf("%d", time.Now().Unix()),
					"10000",
					"2400",
					"604800",
					"3600",
				}, " ")

			rr, err := mdns.NewRR(s)
			if err == nil {
				msg.Answer = append(msg.Answer, rr)
			}

		case "NS":
			for _, server := range dns.Config.DNS.Servers {
				s := strings.Join([]string{
					q.Name,
					"3600",
					"IN",
					"NS",
					server.Name,
				}, " ")

				rr, err := mdns.NewRR(s)
				if err == nil {
					msg.Answer = append(msg.Answer, rr)
				}
			}

			for _, server := range dns.Config.DNS.Servers {
				s := strings.Join([]string{
					server.Name,
					"3600",
					"IN",
					"A",
					server.IP,
				}, " ")
				rr, err := mdns.NewRR(s)
				if err == nil {
					msg.Extra = append(msg.Extra, rr)
				}
			}
		}
	}
	return
}
