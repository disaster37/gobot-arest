package extra

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

const valueTestDelay = 250

var _ gobot.Driver = (*ValueDriver)(nil)

func initTestValueDriver() (*ValueDriver, *extraTestAdaptor) {
	a := newExtraTestAdaptor()

	return NewValueDriver(a, "test", 10*time.Millisecond), a
}

func TestValueDriverDefaultName(t *testing.T) {
	g, _ := initTestValueDriver()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "Value"))
}

func TestValueDriverSetName(t *testing.T) {
	g, _ := initTestValueDriver()
	g.SetName("mybot")
	assert.Equal(t, g.Name(), "mybot")
}

func TestValueDriverValueName(t *testing.T) {
	g, _ := initTestValueDriver()
	assert.Equal(t, g.ValueName(), "test")
}

func TestValueDriverStart(t *testing.T) {
	sem := make(chan bool)
	d, a := initTestValueDriver()

	// Test Read value and wait event
	if err := d.Once(NewValue, func(data interface{}) {
		assert.Equal(t, d.data, 10)
		sem <- true
	}); err != nil {
		t.Fatal(err)
	}
	a.TestAdaptorValueRead(func(name string) (val interface{}, err error) {
		val = 10
		return
	})
	assert.NoError(t, d.Start())
	select {
	case <-sem:
	case <-time.After(valueTestDelay * time.Millisecond):
		t.Errorf("Extra Event \"NewValue\" was not published")
	}

	// Test Read value when error and wait event
	if err := d.Once(Error, func(data interface{}) {
		sem <- true
	}); err != nil {
		t.Fatal(err)
	}
	a.TestAdaptorValueRead(func(name string) (val interface{}, err error) {
		err = errors.New("value read error")
		return
	})
	select {
	case <-sem:
	case <-time.After(valueTestDelay * time.Millisecond):
		t.Errorf("Extra Event \"Error\" was not published")
	}

}

func TestValueDriverHalt(t *testing.T) {
	d, _ := initTestValueDriver()
	go func() {
		<-d.halt
	}()
	assert.NoError(t, d.Halt())
}
