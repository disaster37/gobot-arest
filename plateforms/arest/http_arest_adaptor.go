package arest

import (
	"time"

	restClient "github.com/disaster37/gobot-arest/v1/plateforms/arest/client/rest"
	"gobot.io/x/gobot"
)

// HTTPAdaptor is the Gobot Adaptor for Arest based boards
type HTTPAdaptor struct {
	url string
	Adaptor
}

// NewHTTPAdaptor returns a new HTTP Arest Adaptor which optionally accepts:
//
//	string: The board name
//	time.Duration: The timeout for http backend
//	bool: The debug mode
func NewHTTPAdaptor(url string, args ...interface{}) *Adaptor {
	a := &Adaptor{
		name:    gobot.DefaultName("arest"),
		isDebug: false,
		timeout: 0,
		Eventer: gobot.NewEventer(),
	}

	for _, arg := range args {
		switch arg.(type) {
		case string:
			a.name = arg.(string)
		case time.Duration:
			a.timeout = arg.(time.Duration)
		case bool:
			a.isDebug = arg.(bool)
		}
	}

	a.Board = restClient.NewClient(url, a.timeout, a.isDebug)

	return a
}
