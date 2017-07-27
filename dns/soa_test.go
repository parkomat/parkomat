package dns

import (
	mdns "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSOAHandle(t *testing.T) {
	saved := now
	now = func() int64 {
		return 0
	}
	defer func() {
		now = saved
	}()

	msg := &mdns.Msg{}

	question := mdns.Question{
		Name: "test.com",
	}

	soa := NewSOAHandler("kitty")

	err := soa.Handle(msg, nil, question)
	assert.Nil(t, err)

	expectedMsg := &mdns.Msg{}
	rr, err := mdns.NewRR("test.com. 3600 IN SOA kitty admin.test.com 0 10000 2400 604800 3600")
	assert.Nil(t, err)
	expectedMsg.Answer = append(expectedMsg.Answer, rr)

	assert.Exactly(t, msg.Answer, expectedMsg.Answer)
}
