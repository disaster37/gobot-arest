package arest

import (
	"context"

	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	"gobot.io/x/gobot"
)

type mockArestBoard struct {
	disconnectError error
	gobot.Eventer
	pins map[int]*client.Pin
}

func newMockArestBoard() *mockArestBoard {
	m := &mockArestBoard{
		Eventer:         gobot.NewEventer(),
		disconnectError: nil,
		pins:            make(map[int]*client.Pin, 100),
	}

	m.pins[1] = &client.Pin{Value: 1}
	m.pins[15] = &client.Pin{Value: 133}

	return m
}

func (mockArestBoard) Connect(ctx context.Context) error { return nil }
func (m mockArestBoard) Disconnect(ctx context.Context) error {
	return m.disconnectError
}
func (mockArestBoard) Reconnect(ctx context.Context) error { return nil }
func (m mockArestBoard) Pins() map[int]*client.Pin {
	return m.pins
}
func (mockArestBoard) SetPinMode(ctx context.Context, pin int, mode string) (err error) { return }
func (mockArestBoard) DigitalRead(ctx context.Context, pin int) (level int, err error)  { return }
func (mockArestBoard) DigitalWrite(ctx context.Context, pin int, level int) (err error) { return nil }
func (mockArestBoard) ReadValue(ctx context.Context, name string) (value interface{}, err error) {
	return 10, nil
}
func (mockArestBoard) ReadValues(ctx context.Context) (values map[string]interface{}, err error) {
	return map[string]interface{}{
		"test": 10,
	}, nil
}
func (mockArestBoard) CallFunction(ctx context.Context, name string, param string) (resp int, err error) {
	return
}
func (m mockArestBoard) AddPin(name int, pin *client.Pin) {
	m.pins[name] = pin
}

func initTestAdaptor() *Adaptor {
	a := NewHTTPAdaptor("http://localhost")
	a.Board = newMockArestBoard()
	a.Connect()
	return a
}
