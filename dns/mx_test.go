package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/parkomat/parkomat/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMXHandle(t *testing.T) {
	msg := &mdns.Msg{}

	zone := &config.Zone{
		MX: `1 test1.mail.server
10 test2.mail.server`,
	}

	question := mdns.Question{
		Name: "test.com",
	}

	mx := &mxHandler{}

	err := mx.Handle(msg, zone, question)
	assert.Nil(t, err)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.com. 3600 IN MX 1 test1.mail.server")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)
	rr, err = mdns.NewRR("test.com. 3600 IN MX 10 test2.mail.server")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, msg.Answer, expectedMsg.Answer)
}
