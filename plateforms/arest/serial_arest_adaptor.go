package arest

import (
	"context"
	"time"

	restClient "github.com/disaster37/gobot-arest/v1/plateforms/client/rest"
	"gobot.io/x/gobot"
)

// SerialAdaptor is the Gobot Adaptor for Arest based boards
type SerialAdaptor struct {
	name    string
	port    string
	timeout time.Duration
	isDebug bool
	Board   arestBoard
	gobot.Eventer
}

// NewSerialAdaptor returns a new serial Arest Adaptor which optionally accepts:
//
//	string: The board name
//	time.Duration: The timeout for serial response
//	bool: The debug mode
func NewSerialAdaptor(port string, args ...interface{}) ArestAdaptor {
	a := &HTTPAdaptor{
		name:    gobot.DefaultName("arest"),
		url:     port,
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

	a.Board = restClient.NewClient(a.url, a.timeout, a.isDebug)

	return a
}

// Connect init connection throught HTTP to the board
func (a *SerialAdaptor) Connect(ctx context.Context) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.Board.Connect(ctx)
	}
}

// Disconnect close the connection to the Board
func (a *SerialAdaptor) Disconnect(ctx context.Context) (err error) {
	if a.Board != nil {
		return a.Board.Disconnect(ctx)
	}
	return nil
}

// Finalize terminates the Arest connection
func (a *SerialAdaptor) Finalize(ctx context.Context) (err error) {
	return a.Disconnect(ctx)
}

// Reconnect permit to reopen connection to the board
func (a *SerialAdaptor) Reconnect(ctx context.Context) (err error) {
	return a.Board.Reconnect(ctx)
}

// Name returns the Arest Adaptors name
func (a *SerialAdaptor) Name() string {
	return a.name
}

// SetName sets the Arest Adaptors name
func (a *SerialAdaptor) SetName(name string) {
	a.name = name
}
