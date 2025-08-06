package arest

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*SerialAdaptor)(nil)

func initTestSerialAdaptor() *Adaptor {
	a := NewSerialAdaptor("/dev/null")
	return a
}

func TestArestSerialAdaptor(t *testing.T) {

	// With basic parameters
	a := initTestSerialAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "SerialArest"))

	// With all parameters
	a = NewSerialAdaptor("/dev/null", 10*time.Second, "TEST", true)
	assert.Equal(t, "TEST", a.Name())
	assert.Equal(t, 10.*time.Second, a.timeout)
	assert.True(t, a.isDebug)
}
