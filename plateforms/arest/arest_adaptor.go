package arest

import (
	"context"

	"gobot.io/x/gobot"
)

type arestBoard interface {
	// Connect permit to open connection on board
	Connect(ctx context.Context) error

	// Disconnect permit to close connection on board
	Disconnect(ctx context.Context) error

	// Reconnect permit to reopen connection on board
	Reconnect(ctx context.Context) error

	// SetPinMode permit to set pin mode
	SetPinMode(ctx context.Context, pin int, mode string) (err error)

	// DigitalWrite permit to set level on pin
	DigitalWrite(ctx context.Context, pin int, level int) (err error)

	// DigitalRead permit to read level from pin
	DigitalRead(ctx context.Context, pin int) (level int, err error)

	// ReadValue permit to read user variable
	ReadValue(ctx context.Context, name string) (value interface{}, err error)

	// ReadValues permit to read all user variables
	ReadValues(ctx context.Context) (values map[string]interface{}, err error)

	// CallFunction permit to call user function
	CallFunction(ctx context.Context, name string, param string) (resp int, err error)

	gobot.Eventer
}

// ArestAdaptor represent the arest adaptor interface
type ArestAdaptor interface {
	Connect(ctx context.Context) (err error)
	Finalize(ctx context.Context) (err error)
	Reconnect(ctx context.Context) (err error)
	Name() string
	SetName(n string)
	gobot.Eventer
}
