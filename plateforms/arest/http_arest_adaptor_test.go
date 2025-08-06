package arest

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Adaptor = (*HTTPAdaptor)(nil)

func initTestHTTPAdaptor() *Adaptor {
	a := NewHTTPAdaptor("http://localhost:4567")
	return a
}

func TestArestHTTPAdaptor(t *testing.T) {

	// With minimal parameter
	a := initTestHTTPAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "HTTPArest"))

	// With all parameters
	a = NewHTTPAdaptor("http://localhost", 10*time.Second, "TEST", true)
	assert.Equal(t, "TEST", a.Name())
	assert.Equal(t, 10.*time.Second, a.timeout)
	assert.True(t, a.isDebug)
}
