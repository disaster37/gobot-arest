package arest

import (
	"context"
	"strconv"

	"github.com/disaster37/gobot-arest/plateforms/arest/client"
)

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (a *Adaptor) DigitalWrite(pin string, level byte) (err error) {

	p, err := strconv.Atoi(pin)
	l := int(level)
	ctx := context.TODO()
	if err != nil {
		return err
	}

	if a.Board.Pins()[p] == nil {
		err = a.Board.SetPinMode(ctx, p, client.ModeOutput)
		if err != nil {
			return err
		}
	}

	return a.Board.DigitalWrite(ctx, p, l)
}

// DigitalRead retrieves digital value from specified pin.
// Returns -1 if the response from the board has timed out
func (a *Adaptor) DigitalRead(pin string) (val int, err error) {

	p, err := strconv.Atoi(pin)
	if err != nil {
		return val, err
	}
	ctx := context.TODO()
	

	if a.Board.Pins()[p] == nil {
		if err = a.Board.SetPinMode(ctx, p, client.ModeInput); err != nil {
			return val, err
		}
	}

	return a.Board.DigitalRead(ctx, p)
}
