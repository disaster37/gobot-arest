package arest

import (
	"context"
)

// ValueRead can read on value on plateform like aRest
func (a *Adaptor) ValueRead(name string) (val interface{}, err error) {

	ctx := context.TODO()

	return a.Board.ReadValue(ctx, name)
}

// ValuesRead can read all values on plateform like aRest
func (a *Adaptor) ValuesRead() (vals map[string]interface{}, err error) {
	ctx := context.TODO()
	return a.Board.ReadValues(ctx)
}

// FunctionCall can call function on plateform like aRest
func (a *Adaptor) FunctionCall(name string, parameters string) (val int, err error) {
	ctx := context.TODO()
	return a.Board.CallFunction(ctx, name, parameters)
}
