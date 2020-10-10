package arest

import (
	"context"
	"time"

	"github.com/disaster37/gobot-arest/plateforms/arest/client"
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

	// Pins perit to get current pin settings
	Pins() map[int]*client.Pin

	// AddPin permit to add pin setting
	AddPin(name int, pin *client.Pin)

	gobot.Eventer
}

// ArestAdaptor represent the arest adaptor interface
type ArestAdaptor interface {
	Connect() (err error)
	Finalize() (err error)
	Reconnect() (err error)
	Name() string
	SetName(n string)
	gobot.Eventer
}

// Adaptor is a general Arest Adaptor
type Adaptor struct {
	timeout time.Duration
	isDebug bool
	Board   arestBoard
	gobot.Eventer
	name string
}

// Connect init connection throught HTTP to the board
func (a *Adaptor) Connect() (err error) {
	return a.Board.Connect(context.TODO())
}

// Disconnect close the connection to the Board
func (a *Adaptor) Disconnect() (err error) {
	if a.Board != nil {
		return a.Board.Disconnect(context.TODO())
	}
	return nil
}

// Finalize terminates the Arest connection
func (a *Adaptor) Finalize() (err error) {
	return a.Disconnect()
}

// Reconnect permit to reopen connection to the board
func (a *Adaptor) Reconnect() (err error) {
	return a.Board.Reconnect(context.TODO())
}

// Name returns the Arest Adaptors name
func (a *Adaptor) Name() string {
	return a.name
}

// SetName sets the Arest Adaptors name
func (a *Adaptor) SetName(name string) {
	a.name = name
}
