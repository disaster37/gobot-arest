package serialClient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"gobot.io/x/gobot"
)

// Client implement arest interface
type Client struct {
	serialPort serial.Port
	serialMode *serial.Mode
	isDebug    bool
	port       string
	timeout    time.Duration
	mutex      sync.Mutex
	mutexPin   sync.Mutex

	// It permit to exchange error, result and start watchdog between read routine and write
	com *Com

	// It permit to know if current connexion is connected
	connected atomic.Value
	// It permit to set the right mode with digital read / write
	pins atomic.Value
	gobot.Eventer
}

// NewClient permit to initialize new client Object
func NewClient(port string, serialMode *serial.Mode, timeout time.Duration, isDebug bool) *Client {

	clientArest := &Client{
		serialPort: nil,
		isDebug:    isDebug,
		port:       port,
		serialMode: serialMode,
		timeout:    timeout,
		mutex:      sync.Mutex{},
		mutexPin:   sync.Mutex{},
		connected:  atomic.Value{},
		com: &Com{
			Res:      make(chan string),
			Err:      make(chan error),
			Watchdog: make(chan bool),
		},
		Eventer: gobot.NewEventer(),
		pins:    atomic.Value{},
	}

	clientArest.pins.Store(make(map[int]*client.Pin))

	clientArest.AddEvent("connected")
	clientArest.AddEvent("disconnected")
	clientArest.AddEvent("reconnected")
	clientArest.AddEvent("timeout")
	clientArest.connected.Store(false)

	// It permit to try to reconnect on serial if timeout throw from watchdog
	// It try for ever to reconnect on board
	clientArest.On("timeout", func(s interface{}) {
		isReconnected := false
		for !isReconnected {
			time.Sleep(1 * time.Millisecond)
			err := clientArest.Reconnect(context.TODO())
			if err == nil {
				isReconnected = true
			} else {
				log.Error(err)
			}
		}
	})

	return clientArest
}

// Client permit to get curent serial client
func (c *Client) Client() serial.Port {
	return c.serialPort
}

// Pins return the current pins
func (c *Client) Pins() map[int]*client.Pin {
	return c.pins.Load().(map[int]*client.Pin)
}

// AddPin permit to add pin.
func (c *Client) AddPin(name int, pin *client.Pin) {
	c.mutexPin.Lock()
	defer c.mutexPin.Unlock()

	pins := c.Pins()
	pins[name] = pin
	c.pins.Store(pins)
}

// Connect start connection to the board
// It call the root url to check if board is online
func (c *Client) Connect(ctx context.Context) (err error) {

	if c.connected.Load().(bool) {
		return
	}

	serialPort, err := serial.Open(c.port, c.serialMode)
	if err != nil {
		return err
	}
	c.serialPort = serialPort

	// clean current serial
	serialPort.ResetInputBuffer()
	serialPort.ResetOutputBuffer()

	// Start routine to read serial
	c.readProcess(ctx)

	// Try connexion
	time.Sleep(1 * time.Second)
	_, err = c.write(ctx, "/")
	if err != nil {
		return err
	}

	c.Publish("connected", true)
	c.connected.Store(true)

	return nil
}

// Disconnect close connecion to the board
func (c *Client) Disconnect(ctx context.Context) (err error) {

	c.connected.Store(false)
	err = c.serialPort.Close()
	c.serialPort.ResetInputBuffer()
	c.serialPort.ResetOutputBuffer()

	if err != nil {
		return err
	}

	c.Publish("disconnected", true)

	return nil
}

// Reconnect close and start connection to the board
func (c *Client) Reconnect(ctx context.Context) (err error) {
	err = c.Disconnect(ctx)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	err = c.Connect(ctx)
	if err != nil {
		return err
	}

	// Set pin mode and output
	for pin, state := range c.Pins() {
		err = c.SetPinMode(ctx, pin, state.Mode)
		if err != nil {
			return err
		}

		if state.Mode == client.ModeOutput {
			err = c.DigitalWrite(ctx, pin, state.Value)
			if err != nil {
				return err
			}
		}
	}

	c.Publish("reconnected", true)

	return nil
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(ctx context.Context, pin int, mode string) (err error) {

	if !c.connected.Load().(bool) {
		return errors.New("Not connected")
	}

	if c.Pins()[pin] == nil {
		c.AddPin(pin, &client.Pin{})
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Pin: %d, Mode: %s", pin, mode)
		}

		if mode != client.ModeInput && mode != client.ModeInputPullup && mode != client.ModeOutput {
			return errors.Errorf("Can't found mode %s", mode)
		}

		url := fmt.Sprintf("/mode/%d/%s\n\r", pin, mode)

		resp, err := c.write(ctx, url)
		if err != nil {
			return err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		c.Pins()[pin].Mode = mode

		return nil
	}
}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(ctx context.Context, pin int, level int) (err error) {

	if !c.connected.Load().(bool) {
		return errors.New("Not connected")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Pin: %d, Level: %d", pin, level)
		}

		if level != client.LevelHigh && level != client.LevelLow {
			return errors.Errorf("Can't found level %d", level)
		}

		url := fmt.Sprintf("/digital/%d/%d\n\r", pin, level)

		resp, err := c.write(ctx, url)
		if err != nil {
			return err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		c.Pins()[pin].Value = level

		return nil
	}
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(ctx context.Context, pin int) (level int, err error) {

	if !c.connected.Load().(bool) {
		return level, errors.New("Not connected")
	}

	select {
	case <-ctx.Done():
		return level, ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Pin: %d", pin)
		}

		url := fmt.Sprintf("/digital/%d\n\r", pin)
		data := make(map[string]interface{})

		resp, err := c.write(ctx, url)
		if err != nil {
			return level, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return level, err
		}

		return int(data["return_value"].(float64)), nil
	}
}

// ReadValue permit to read user variable
func (c *Client) ReadValue(ctx context.Context, name string) (value interface{}, err error) {
	if !c.connected.Load().(bool) {
		return value, errors.New("Not connected")
	}

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Value name: %s", name)
		}

		url := fmt.Sprintf("/%s\n\r", name)
		data := make(map[string]interface{})

		resp, err := c.write(ctx, url)
		if err != nil {
			return nil, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return nil, err
		}

		if temp, ok := data[name]; ok {
			value = temp
		} else {
			err = errors.Errorf("Variable %s not found", name)
		}

		return value, err
	}
}

// ReadValues permit to read user variable
func (c *Client) ReadValues(ctx context.Context) (values map[string]interface{}, err error) {
	if !c.connected.Load().(bool) {
		return values, errors.New("Not connected")
	}

	select {
	case <-ctx.Done():
		return values, ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()
		url := "/\n\r"
		data := make(map[string]interface{})

		resp, err := c.write(ctx, url)
		if err != nil {
			return nil, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return nil, err
		}

		if temp, ok := data["variables"]; ok {
			values = temp.(map[string]interface{})
		} else {
			err = errors.Errorf("No variable found")
		}

		return values, err
	}
}

// CallFunction permit to call user function
func (c *Client) CallFunction(ctx context.Context, name string, param string) (value int, err error) {
	if !c.connected.Load().(bool) {
		return value, errors.New("Not connected")
	}

	select {
	case <-ctx.Done():
		return value, ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Function: %s, param: %s", name, param)
		}

		url := fmt.Sprintf("/%s?params=%s\n\r", name, param)
		data := make(map[string]interface{})

		resp, err := c.write(ctx, url)
		if err != nil {
			return value, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		err = json.Unmarshal([]byte(resp), &data)
		if err != nil {
			return value, err
		}

		if temp, ok := data["return_value"]; ok {
			value = int(temp.(float64))
		} else {
			err = errors.Errorf("Function %s not found", name)
		}

		return value, nil
	}
}
