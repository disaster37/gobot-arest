package extra

import (
	"errors"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

const valuesTestDelay = 250

var _ gobot.Driver = (*ValuesDriver)(nil)

func initTestValuesDriver() (*ValuesDriver, *extraTestAdaptor) {
	a := newExtraTestAdaptor()

	return NewValuesDriver(a, 10*time.Millisecond), a
}

func TestValuesDriverDefaultName(t *testing.T) {
	g, _ := initTestValuesDriver()
	gobottest.Refute(t, g.Connection(), nil)
	gobottest.Assert(t, strings.HasPrefix(g.Name(), "Values"), true)
}

func TestValuesDriverSetName(t *testing.T) {
	g, _ := initTestValuesDriver()
	g.SetName("mybot")
	gobottest.Assert(t, g.Name(), "mybot")
}

func TestValuesDriverStart(t *testing.T) {
	sem := make(chan bool, 0)
	d, a := initTestValuesDriver()

	// Test Read value and wait event
	d.Once(NewValues, func(data interface{}) {
		gobottest.Assert(t, d.data, map[string]interface{}{
			"test": 10,
		})
		sem <- true
	})
	a.TestAdaptorValuesRead(func() (vals map[string]interface{}, err error) {
		vals = map[string]interface{}{
			"test": 10,
		}
		return
	})
	gobottest.Assert(t, d.Start(), nil)
	select {
	case <-sem:
	case <-time.After(valuesTestDelay * time.Millisecond):
		t.Errorf("Extra Event \"NewValues\" was not published")
	}

	// Test Read values when error and wait event
	d.Once(Error, func(data interface{}) {
		sem <- true
	})
	a.TestAdaptorValuesRead(func() (vals map[string]interface{}, err error) {
		err = errors.New("values read error")
		return
	})
	select {
	case <-sem:
	case <-time.After(valuesTestDelay * time.Millisecond):
		t.Errorf("Extra Event \"Error\" was not published")
	}
}

func TestValuesDriverHalt(t *testing.T) {
	d, _ := initTestValuesDriver()
	go func() {
		<-d.halt
	}()
	gobottest.Assert(t, d.Halt(), nil)
}
