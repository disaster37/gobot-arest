package arest

import (
	"context"
	"time"

	restClient "github.com/disaster37/gobot-arest/v1/plateforms/client/rest"
	"gobot.io/x/gobot"
)

// HTTPAdaptor is the Gobot Adaptor for Arest based boards
type HTTPAdaptor struct {
	name    string
	url     string
	timeout time.Duration
	isDebug bool
	Board   arestBoard
	gobot.Eventer
}

// NewHTTPAdaptor returns a new HTTP Arest Adaptor which optionally accepts:
//
//	string: The board name
//	time.Duration: The timeout for http backend
//	bool: The debug mode
func NewHTTPAdaptor(url string, args ...interface{}) ArestAdaptor {
	a := &HTTPAdaptor{
		name:    gobot.DefaultName("arest"),
		url:     url,
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
func (a *HTTPAdaptor) Connect(ctx context.Context) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return a.Board.Connect(ctx)
	}
}

// Disconnect close the connection to the Board
func (a *HTTPAdaptor) Disconnect(ctx context.Context) (err error) {
	if a.Board != nil {
		return a.Board.Disconnect(ctx)
	}
	return nil
}

// Finalize terminates the Arest connection
func (a *HTTPAdaptor) Finalize(ctx context.Context) (err error) {
	return a.Disconnect(ctx)
}

// Reconnect permit to reopen connection to the board
func (a *HTTPAdaptor) Reconnect(ctx context.Context) (err error) {
	return a.Board.Reconnect(ctx)
}

// Name returns the Arest Adaptors name
func (a *HTTPAdaptor) Name() string {
	return a.name
}

// SetName sets the Arest Adaptors name
func (a *HTTPAdaptor) SetName(name string) {
	a.name = name
}
