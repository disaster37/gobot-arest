package arest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ ArestAdaptor = (*Adaptor)(nil)
var _ extra.ExtraReader = (*Adaptor)(nil)

type mockFirmataBoard struct {
	disconnectError error
	gobot.Eventer
	pins map[int]*client.Pin
}

func newMockFirmataBoard() *mockFirmataBoard {
	m := &mockFirmataBoard{
		Eventer:         gobot.NewEventer(),
		disconnectError: nil,
		pins:            make(map[int]*client.Pin, 100),
	}

	m.pins[1] = &client.Pin{Value: 1}
	m.pins[15] = &client.Pin{Value: 133}

	return m
}

func (mockFirmataBoard) Connect(ctx context.Context) error { return nil }
func (m mockFirmataBoard) Disconnect(ctx context.Context) error {
	return m.disconnectError
}
func (mockFirmataBoard) Reconnect(ctx context.Context) error { return nil }
func (m mockFirmataBoard) Pins() map[int]*client.Pin {
	return m.pins
}
func (mockFirmataBoard) SetPinMode(ctx context.Context, pin int, mode string) (err error) { return }
func (mockFirmataBoard) DigitalRead(ctx context.Context, pin int) (level int, err error)  { return }
func (mockFirmataBoard) DigitalWrite(ctx context.Context, pin int, level int) (err error) { return nil }
func (mockFirmataBoard) ReadValue(ctx context.Context, name string) (value interface{}, err error) {
	return 10, nil
}
func (mockFirmataBoard) ReadValues(ctx context.Context) (values map[string]interface{}, err error) {
	return map[string]interface{}{
		"test": 10,
	}, nil
}
func (mockFirmataBoard) CallFunction(ctx context.Context, name string, param string) (resp int, err error) {
	return
}
func (m mockFirmataBoard) AddPin(name int, pin *client.Pin) {
	m.pins[name] = pin
}

func initTestAdaptor() *Adaptor {
	a := NewHTTPAdaptor("http://localhost")
	a.Board = newMockFirmataBoard()
	a.Connect()
	return a
}

func TestAdaptor(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "HTTPArest"), true)
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)

	a = initTestAdaptor()
	a.Board.(*mockFirmataBoard).disconnectError = errors.New("close error")
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestAdaptorConnect(t *testing.T) {

	// Without error
	a := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)

	// Disconnect
	a = initTestAdaptor()
	gobottest.Assert(t, a.Disconnect(), nil)

	// Reconnect
	a = initTestAdaptor()
	gobottest.Assert(t, a.Reconnect(), nil)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.DigitalWrite("1", 1), nil)
}

func TestAdaptorDigitalWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Refute(t, a.DigitalWrite("xyz", 50), nil)
}

func TestAdaptorDigitalRead(t *testing.T) {
	a := initTestAdaptor()

	val, err := a.DigitalRead("0")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, val, 0)
}

func TestAdaptorDigitalReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.DigitalRead("xyz")
	gobottest.Refute(t, err, nil)
}

func TestAdaptorSetPinMode(t *testing.T) {
	a := initTestAdaptor()

	gobottest.Assert(t, a.Board.SetPinMode(context.Background(), 1, client.ModeInput), nil)
}

func TestAdaptorValueRead(t *testing.T) {
	a := initTestAdaptor()

	value, err := a.ValueRead("test")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, value, 10)
}

func TestAdaptorValuesRead(t *testing.T) {
	a := initTestAdaptor()
	expected := map[string]interface{}{
		"test": 10,
	}
	values, err := a.ValuesRead()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, values, expected)
}

func TestAdaptorFunctionCall(t *testing.T) {
	a := initTestAdaptor()

	value, err := a.FunctionCall("test", "param1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, value, 0)
}
