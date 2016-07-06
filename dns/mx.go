package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"strings"
)

type mxHandler struct {
}

// Handle produces reply for MX question
func (m *mxHandler) Handle(msg *mdns.Msg, zone *config.Zone, question mdns.Question) (err error) {
	for _, s := range strings.Split(zone.MX, "\n") {
		s = strings.Trim(s, " ")
		if s != "" {
			mx := strings.Split(s, " ")

			r := strings.Join([]string{
				question.Name,
				"3600",
				"IN",
				"MX",
				mx[0],
				mx[1],
			}, " ")

			rr, err := mdns.NewRR(r)
			if err == nil {
				msg.Answer = append(msg.Answer, rr)
			}
		}
	}
	return
}
