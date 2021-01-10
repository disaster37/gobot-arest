package arest

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*SerialAdaptor)(nil)

func initTestSerialAdaptor() *Adaptor {
	a := NewSerialAdaptor("/dev/null")
	return a
}

func TestArestSerialAdaptor(t *testing.T) {

	// With basic parameters
	a := initTestSerialAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "SerialArest"), true)

	// With all parameters
	a = NewSerialAdaptor("/dev/null", 10*time.Second, "TEST", true)
	gobottest.Assert(t, "TEST", a.Name())
	gobottest.Assert(t, 10.*time.Second, a.timeout)
	gobottest.Assert(t, true, a.isDebug)
}
