package arest

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*HTTPAdaptor)(nil)

func initTestHTTPAdaptor() *Adaptor {
	a := NewHTTPAdaptor("http://localhost:4567")
	return a
}

func TestArestHTTPAdaptor(t *testing.T) {

	// With minimal parameter
	a := initTestHTTPAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "HTTPArest"), true)

	// With all parameters
	a = NewHTTPAdaptor("http://localhost", 10*time.Second, "TEST", true)
	gobottest.Assert(t, "TEST", a.Name())
	gobottest.Assert(t, 10.*time.Second, a.timeout)
	gobottest.Assert(t, true, a.isDebug)
}
