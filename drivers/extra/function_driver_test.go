package extra

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*FunctionDriver)(nil)

func initTestFunctionDriver() (*FunctionDriver, *extraTestAdaptor) {
	a := newExtraTestAdaptor()

	return NewFunctionDriver(a, "test", "param1"), a
}

func TestFunctionDriverDefaultName(t *testing.T) {
	g, _ := initTestFunctionDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Function"), true)
}

func TestFunctionDriverSetName(t *testing.T) {
	g, _ := initTestFunctionDriver()
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}

func TestFunctionDriverFunctionName(t *testing.T) {
	g, _ := initTestFunctionDriver()
	gobottest.Assert(t, g.FunctionName(), "test")
}

func TestFunctionDriverParameters(t *testing.T) {
	g, _ := initTestFunctionDriver()
	gobottest.Assert(t, g.Parameters(), "param1")
	g.SetParameters("param2")
	gobottest.Assert(t, g.Parameters(), "param2")
}

func TestFunctionDriverStart(t *testing.T) {
	d, _ := initTestFunctionDriver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestFunctionDriverHalt(t *testing.T) {
	d, _ := initTestFunctionDriver()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestFunctionDriverCall(t *testing.T) {
	d, a := initTestFunctionDriver()

	// When return code match with expected
	gobottest.Assert(t, d.Call(), nil)

	// When return code mispatch expected
	a.TestAdaptorFunctionCall(func(name string, parameters string) (val int, err error) {
		return 1, nil
	})
	gobottest.Assert(t, d.Call(), ErrCallFunctionReturnCodeMismatch)

	// When other kind of error
	err := errors.New("test")
	a.TestAdaptorFunctionCall(func(name string, parameters string) (int, error) {
		return 0, err
	})
	gobottest.Assert(t, d.Call(), err)

}
