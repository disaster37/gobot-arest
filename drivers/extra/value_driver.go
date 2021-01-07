package extra

import (
	"time"

	"gobot.io/x/gobot"
)

const (
	// Error event
	Error = "error"

	NewValue = "new-value"
)

type ValueDriver struct {
	valueName  string
	data       interface{}
	name       string
	halt       chan bool
	interval   time.Duration
	connection ExtraValueReader
	gobot.Eventer
}

// NewValueDriver returns a new ValueDriver with a polling interval of
// 1 seconds given a ValueReader.
//
// Optionally accepts:
//  time.Duration: Interval at which the ValueDriver is polled for new information
func NewValueDriver(a ExtraValueReader, valueName string, v ...time.Duration) *ValueDriver {
	b := &ValueDriver{
		name:       gobot.DefaultName("Value"),
		connection: a,
		valueName:  valueName,
		Eventer:    gobot.NewEventer(),
		interval:   10 * time.Second,
		halt:       make(chan bool),
	}

	if len(v) > 0 {
		b.interval = v[0]
	}

	b.AddEvent(NewValue)
	b.AddEvent(Error)

	return b
}

// Start starts the ButtonDriver and polls the state of the button at the given interval.
//
// Emits the Events:
// 	NewValue interface{} - The new value
//	Error error - On button error
func (b *ValueDriver) Start() (err error) {
	go func() {
		for {
			newValue, err := b.connection.ValueRead(b.valueName)
			if err != nil {
				b.Publish(Error, err)
			} else if newValue != b.data {
				b.data = newValue
				b.Publish(NewValue, newValue)
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

// Halt stops polling the button for new information
func (b *ValueDriver) Halt() (err error) {
	b.halt <- true
	return
}

// Name returns the ButtonDrivers name
func (b *ValueDriver) Name() string { return b.name }

// SetName sets the ButtonDrivers name
func (b *ValueDriver) SetName(n string) { b.name = n }

// ValueName returns the ValueDriver name
func (b *ValueDriver) ValueName() string { return b.valueName }

// Connection returns the ButtonDrivers Connection
func (b *ValueDriver) Connection() gobot.Connection { return b.connection.(gobot.Connection) }
