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

func TestAdaptor(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "HTTPArest"), true)
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Finalize(), nil)

	a = initTestAdaptor()
	a.Board.(*mockArestBoard).disconnectError = errors.New("close error")
	gobottest.Assert(t, a.Finalize(), errors.New("close error"))
}

func TestAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	a.SetName("test")
	gobottest.Assert(t, "test", a.Name())
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
