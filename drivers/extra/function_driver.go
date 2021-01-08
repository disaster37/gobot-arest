package extra

import (
	"gobot.io/x/gobot"
)

// FunctionDriver represent a function driver
type FunctionDriver struct {
	functionName string
	parameters   string
	returnCode   int
	name         string
	connection   ExtraReader
	gobot.Eventer
}

// NewFunctionDriver returns a new FunctionDriver can be call on demand
func NewFunctionDriver(a ExtraReader, functionName string, parameters string) *FunctionDriver {
	b := &FunctionDriver{
		name:         gobot.DefaultName("Function"),
		connection:   a,
		functionName: functionName,
		parameters:   parameters,
		Eventer:      gobot.NewEventer(),
	}

	b.AddEvent(Error)

	return b
}

// Start implements the Driver interface
//
// Emits the Events:
//	Error error - On function call error
func (b *FunctionDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (b *FunctionDriver) Halt() (err error) { return }

// Name returns the FunctionDriver name
func (b *FunctionDriver) Name() string { return b.name }

// SetName sets the FunctionDriver name
func (b *FunctionDriver) SetName(n string) { b.name = n }

// FunctionName returns the FunctionDriver name
func (b *FunctionDriver) FunctionName() string { return b.functionName }

// Parameters returns the FunctionDriver parameters
func (b *FunctionDriver) Parameters() string { return b.parameters }

// SetParameters set the FunctionDriver parameters
func (b *FunctionDriver) SetParameters(parameters string) { b.parameters = parameters }

// Connection returns the FunctionDriver Connection
func (b *FunctionDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }

// Call run function
//
// Emits the Events:
//	Error error - On function call error
func (b *FunctionDriver) Call() (err error) {
	ret, err := b.connection.FunctionCall(b.functionName, b.parameters)
	if err != nil {
		b.Publish(Error, err)
		return err
	}

	if ret != b.returnCode {
		b.Publish(Error, ErrCallFunctionReturnCodeMismatch)
		return ErrCallFunctionReturnCodeMismatch
	}

	return
}
