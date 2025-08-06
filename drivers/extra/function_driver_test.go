package extra

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*FunctionDriver)(nil)

func initTestFunctionDriver() (*FunctionDriver, *extraTestAdaptor) {
	a := newExtraTestAdaptor()

	return NewFunctionDriver(a, "test", "param1"), a
}

func TestFunctionDriverDefaultName(t *testing.T) {
	g, _ := initTestFunctionDriver()

	g.Connection()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "Function"))
}

func TestFunctionDriverSetName(t *testing.T) {
	g, _ := initTestFunctionDriver()
	g.SetName("mybot")
	assert.Equal(t, g.Name(), "mybot")
}

func TestFunctionDriverFunctionName(t *testing.T) {
	g, _ := initTestFunctionDriver()
	assert.Equal(t, g.FunctionName(), "test")
}

func TestFunctionDriverParameters(t *testing.T) {
	g, _ := initTestFunctionDriver()
	assert.Equal(t, g.Parameters(), "param1")
	g.SetParameters("param2")
	assert.Equal(t, g.Parameters(), "param2")
}

func TestFunctionDriverStart(t *testing.T) {
	d, _ := initTestFunctionDriver()
	assert.NoError(t, d.Start())
}

func TestFunctionDriverHalt(t *testing.T) {
	d, _ := initTestFunctionDriver()
	assert.NoError(t, d.Halt())
}

func TestFunctionDriverCall(t *testing.T) {
	d, a := initTestFunctionDriver()

	// When return code match with expected
	assert.NoError(t, d.Call())

	// When return code mispatch expected
	a.TestAdaptorFunctionCall(func(name string, parameters string) (val int, err error) {
		return 1, nil
	})
	assert.Error(t, d.Call(), ErrCallFunctionReturnCodeMismatch)

	// When other kind of error
	err := errors.New("test")
	a.TestAdaptorFunctionCall(func(name string, parameters string) (int, error) {
		return 0, err
	})
	assert.Error(t, d.Call(), err)

}
