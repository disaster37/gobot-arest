package arest

import (
	"time"

	restClient "github.com/disaster37/gobot-arest/plateforms/arest/client/rest"
	"gobot.io/x/gobot"
)

// HTTPAdaptor is the Gobot Adaptor for Arest based boards
type HTTPAdaptor struct {
	Adaptor
}

// NewHTTPAdaptor returns a new HTTP Arest Adaptor which optionally accepts:
//
//	string: The board name
//	time.Duration: The timeout for http backend
//	bool: The debug mode
func NewHTTPAdaptor(url string, args ...interface{}) *Adaptor {
	a := &Adaptor{
		name:    gobot.DefaultName("HTTPArest"),
		isDebug: false,
		timeout: 0,
		Eventer: gobot.NewEventer(),
	}

	for _, arg := range args {
		switch argTmp := arg.(type) {
		case string:
			a.name = argTmp
		case time.Duration:
			a.timeout = argTmp
		case bool:
			a.isDebug = argTmp
		}
	}

	a.Board = restClient.NewClient(url, a.timeout, a.isDebug)

	return a
}
