package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNSHandle(t *testing.T) {
	msg := &mdns.Msg{}

	question := mdns.Question{
		Name: "test.com",
	}

	c := &config.Config{
		DNS: config.DNS{
			Servers: []config.Server{
				{
					Name: "ns1.test.com",
					IP:   "127.0.0.1",
				},
				{
					Name: "ns2.test.com",
					IP:   "127.0.0.2",
				},
			},
		},
	}

	ns := NewNSHandler(c)

	err := ns.Handle(msg, nil, question)
	assert.Nil(t, err)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.com. 3600 IN NS ns1.test.com")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)
	rr, err = mdns.NewRR("test.com. 3600 IN NS ns2.test.com")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, msg.Answer, expectedMsg.Answer)

	rr, err = mdns.NewRR("ns1.test.com. 3600 IN A 127.0.0.1")
	assert.Nil(t, err)
	expectedMsg.Extra = append(expectedMsg.Extra, rr)
	rr, err = mdns.NewRR("ns2.test.com. 3600 IN A 127.0.0.2")
	assert.Nil(t, err)
	expectedMsg.Extra = append(expectedMsg.Extra, rr)

	assert.Exactly(t, msg.Extra, expectedMsg.Extra)
}
