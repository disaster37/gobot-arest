package arest

import (
	"context"
)

// ValueRead can read on value on plateform
func (a *Adaptor) ValueRead(name string) (val interface{}, err error) {

	ctx := context.TODO()

	return a.Board.ReadValue(ctx, name)
}
