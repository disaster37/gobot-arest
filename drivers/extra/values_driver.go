package extra

import (
	"reflect"
	"time"

	"gobot.io/x/gobot"
)

// ValuesDriver represent values driver
type ValuesDriver struct {
	data       map[string]interface{}
	name       string
	halt       chan bool
	interval   time.Duration
	connection ExtraReader
	gobot.Eventer
}

// NewValuesDriver returns a new ValuesDriver with a polling interval of
// 1 seconds given a ValuesReader.
//
// Optionally accepts:
//  time.Duration: Interval at which the ValuesDriver is polled for new information
func NewValuesDriver(a ExtraReader, v ...time.Duration) *ValuesDriver {
	b := &ValuesDriver{
		name:       gobot.DefaultName("Values"),
		connection: a,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Second,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(NewValues)
	b.AddEvent(Error)

	return b
}

// Start starts the ValuesDriver and polls the new values at the given interval.
//
// Emits the Events:
// 	NewValues map[string]interface{} - The new values
//	Error error - On values error
func (b *ValuesDriver) Start() (err error) {
	go func() {
		for {
			newValues, err := b.connection.ValuesRead()
			if err != nil {
				b.Publish(Error, err)
			} else if !reflect.DeepEqual(b.data, newValues) {
				b.data = newValues
				b.Publish(NewValues, newValues)
			}
			select {
			case <-time.After(b.interval):
			case <-b.halt:
				return
			}
		}
	}()
	return
}

// Halt stops polling the values for new information
func (b *ValuesDriver) Halt() (err error) {
	b.halt <- true
	return
}

// Name returns the ValuesDrivers name
func (b *ValuesDriver) Name() string { return b.name }

// SetName sets the ValuesDrivers name
func (b *ValuesDriver) SetName(n string) { b.name = n }

// Connection returns the ValuesDrivers Connection
func (b *ValuesDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }
