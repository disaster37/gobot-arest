package arest

import (
	"time"

	serialClient "github.com/disaster37/gobot-arest/plateforms/arest/client/serial"
	"go.bug.st/serial"
	"gobot.io/x/gobot"
)

// SerialAdaptor is the Gobot Adaptor for Arest based boards
type SerialAdaptor struct {
	port string
	Adaptor
}

// NewSerialAdaptor returns a new serial Arest Adaptor which optionally accepts:
//
//	string: The board name
//	time.Duration: The timeout for serial response
//	bool: The debug mode
// serial.Mode: the serial mode
func NewSerialAdaptor(port string, args ...interface{}) *Adaptor {
	a := &Adaptor{
		name:    gobot.DefaultName("arest"),
		isDebug: false,
		timeout: 0,
		Eventer: gobot.NewEventer(),
	}

	mode := serial.Mode{
		BaudRate: 115200,
	}

	for _, arg := range args {
		switch arg.(type) {
		case string:
			a.name = arg.(string)
		case time.Duration:
			a.timeout = arg.(time.Duration)
		case bool:
			a.isDebug = arg.(bool)
		case serial.Mode:
			mode = arg.(serial.Mode)
		}
	}

	a.Board = serialClient.NewClient(port, &mode, a.timeout, a.isDebug)

	return a
}
