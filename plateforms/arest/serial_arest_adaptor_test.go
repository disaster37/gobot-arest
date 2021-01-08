package arest

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*SerialAdaptor)(nil)

func initTestSerialAdaptor() *Adaptor {
	a := NewSerialAdaptor("/dev/tty1")
	return a
}

func TestArestSerialAdaptor(t *testing.T) {
	a := initTestSerialAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "SerialArest"), true)
}
