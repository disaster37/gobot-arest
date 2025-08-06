package arest

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ ArestAdaptor = (*Adaptor)(nil)
var _ extra.ExtraReader = (*Adaptor)(nil)

func TestAdaptor(t *testing.T) {
	a := initTestAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "HTTPArest"))
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	assert.NoError(t, a.Finalize())

	a = initTestAdaptor()
	a.Board.(*mockArestBoard).disconnectError = errors.New("close error")
	assert.Error(t, a.Finalize(), errors.New("close error"))
}

func TestAdaptorName(t *testing.T) {
	a := initTestAdaptor()
	a.SetName("test")
	assert.Equal(t, "test", a.Name())
}

func TestAdaptorConnect(t *testing.T) {

	// Without error
	a := initTestAdaptor()
	assert.NoError(t, a.Connect())

	// Disconnect
	a = initTestAdaptor()
	assert.NoError(t, a.Disconnect())

	// Reconnect
	a = initTestAdaptor()
	assert.NoError(t, a.Reconnect())
}

func TestAdaptorDigitalWrite(t *testing.T) {
	a := initTestAdaptor()
	assert.NoError(t, a.DigitalWrite("1", 1))
}

func TestAdaptorDigitalWriteBadPin(t *testing.T) {
	a := initTestAdaptor()
	assert.Error(t, a.DigitalWrite("xyz", 50))
}

func TestAdaptorDigitalRead(t *testing.T) {
	a := initTestAdaptor()

	val, err := a.DigitalRead("0")
	assert.NoError(t, err)
	assert.Equal(t, val, 0)
}

func TestAdaptorDigitalReadBadPin(t *testing.T) {
	a := initTestAdaptor()
	_, err := a.DigitalRead("xyz")
	assert.Error(t, err)
}

func TestAdaptorSetPinMode(t *testing.T) {
	a := initTestAdaptor()

	assert.NoError(t, a.Board.SetPinMode(context.Background(), 1, client.ModeInput))
}

func TestAdaptorValueRead(t *testing.T) {
	a := initTestAdaptor()

	value, err := a.ValueRead("test")
	assert.NoError(t, err)
	assert.Equal(t, value, 10)
}

func TestAdaptorValuesRead(t *testing.T) {
	a := initTestAdaptor()
	expected := map[string]interface{}{
		"test": 10,
	}
	values, err := a.ValuesRead()
	assert.NoError(t, err)
	assert.Equal(t, values, expected)
}

func TestAdaptorFunctionCall(t *testing.T) {
	a := initTestAdaptor()

	value, err := a.FunctionCall("test", "param1")
	assert.NoError(t, err)
	assert.Equal(t, value, 0)
}
