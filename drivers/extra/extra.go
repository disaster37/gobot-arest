package extra

import "errors"

const (
	// Error event
	Error = "error"

	// NewValue event
	NewValue = "new-value"

	// NewValues event
	NewValues = "new-values"
)

// ErrCallFunctionReturnCodeMismatch is error when call function return code that mistmach the provided
var ErrCallFunctionReturnCodeMismatch error = errors.New("Return code mismatch when call function")

// ExtraReader can read abitrary value
type ExtraReader interface {
	ValueRead(name string) (val interface{}, err error)
	ValuesRead() (vals map[string]interface{}, err error)
	FunctionCall(name string, parameters string) (val int, err error)
}
