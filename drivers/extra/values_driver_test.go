package extra

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

const valuesTestDelay = 250

var _ gobot.Driver = (*ValuesDriver)(nil)

func initTestValuesDriver() (*ValuesDriver, *extraTestAdaptor) {
	a := newExtraTestAdaptor()

	return NewValuesDriver(a, 10*time.Millisecond), a
}

func TestValuesDriverDefaultName(t *testing.T) {
	g, _ := initTestValuesDriver()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "Values"))
}

func TestValuesDriverSetName(t *testing.T) {
	g, _ := initTestValuesDriver()
	g.SetName("mybot")
	assert.Equal(t, g.Name(), "mybot")
}

func TestValuesDriverStart(t *testing.T) {
	sem := make(chan bool)
	d, a := initTestValuesDriver()

	// Test Read value and wait event
	if err := d.Once(NewValues, func(data interface{}) {
		assert.Equal(t, d.data, map[string]interface{}{
			"test": 10,
		})
		sem <- true
	}); err != nil {
		t.Fatal(err)
	}
	a.TestAdaptorValuesRead(func() (vals map[string]interface{}, err error) {
		vals = map[string]interface{}{
			"test": 10,
		}
		return
	})
	assert.NoError(t, d.Start())
	select {
	case <-sem:
	case <-time.After(valuesTestDelay * time.Millisecond):
		t.Errorf("Extra Event \"NewValues\" was not published")
	}

	// Test Read values when error and wait event
	if err := d.Once(Error, func(data interface{}) {
		sem <- true
	}); err != nil {
		t.Fatal(err)
	}
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
	assert.NoError(t, d.Halt())
}
