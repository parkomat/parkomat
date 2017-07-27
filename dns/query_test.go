package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryHandleA(t *testing.T) {
	c := &config.Config{
		Domains: []*config.Domain{
			{
				Name: "test.com",
				Zone: &config.Zone{
					A: "127.0.0.1",
				},
			},
		},
	}

	dns := NewDNS(c)

	m := &mdns.Msg{}
	m.SetQuestion("test.com.", mdns.TypeA)

	dns.ParseQuery(m)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.com. 3600 IN A 127.0.0.1")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, m.Answer, expectedMsg.Answer)
}

func TestQueryHandleADNS0x20(t *testing.T) {
	c := &config.Config{
		Domains: []*config.Domain{
			{
				Name: "test.com",
				Zone: &config.Zone{
					A: "127.0.0.1",
				},
			},
		},
	}

	dns := NewDNS(c)

	m := &mdns.Msg{}
	m.SetQuestion("test.COM.", mdns.TypeA)

	dns.ParseQuery(m)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.COM. 3600 IN A 127.0.0.1")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, m.Answer, expectedMsg.Answer)
}
