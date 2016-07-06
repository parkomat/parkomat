package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"strings"
)

type nsHandler struct {
	config *config.Config
}

func NewNSHandler(config *config.Config) *nsHandler {
	return &nsHandler{
		config: config,
	}
}

// Handle produces reply for NS question
func (n *nsHandler) Handle(msg *mdns.Msg, zone *config.Zone, question mdns.Question) (err error) {
	for _, server := range n.config.DNS.Servers {
		s := strings.Join([]string{
			question.Name,
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

	for _, server := range n.config.DNS.Servers {
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
	return
}
