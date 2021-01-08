package arest

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Adaptor = (*HTTPAdaptor)(nil)

func initTestHTTPAdaptor() *Adaptor {
	a := NewHTTPAdaptor("http://localhost:4567")
	return a
}

func TestArestHTTPAdaptor(t *testing.T) {
	a := initTestHTTPAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "HTTPArest"), true)
}
