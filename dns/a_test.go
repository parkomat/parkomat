package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAHandle(t *testing.T) {
	msg := &mdns.Msg{}

	zone := &config.Zone{
		A: "127.0.0.1",
	}

	question := mdns.Question{
		Name: "test.com",
	}

	a := &aHandler{}

	err := a.Handle(msg, zone, question)
	assert.Nil(t, err)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.com. 3600 IN A 127.0.0.1")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, msg.Answer, expectedMsg.Answer)
}
