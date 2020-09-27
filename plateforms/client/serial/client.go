package serialClient

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/disaster37/gobot-arest/v1/plateforms/client"
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
	gobot.Eventer
}

// NewClient permit to initialize new client Object
func NewClient(port string, serialMode *serial.Mode, timeout time.Duration, isDebug bool) *Client {

	client := &Client{
		serialPort: nil,
		isDebug:    isDebug,
		port:       port,
		serialMode: serialMode,
		timeout:    timeout,
		mutex:      sync.Mutex{},
		Eventer:    gobot.NewEventer(),
	}

	client.AddEvent("connected")
	client.AddEvent("disconnected")
	client.AddEvent("reconnected")

	return client
}

// Client permit to get curent serial client
func (c *Client) Client() serial.Port {
	return c.serialPort
}

// Connect start connection to the board
// We just try read port 0 value to check http connexion is ready
func (c *Client) Connect(ctx context.Context) (err error) {

	serialPort, err := serial.Open(c.port, c.serialMode)
	if err != nil {
		return err
	}

	// clean current serial
	serialPort.ResetInputBuffer()
	serialPort.ResetOutputBuffer()

	_, err = c.DigitalRead(ctx, 0)
	if err != nil {
		return err
	}

	c.Publish("connected", true)
	return nil
}

// Disconnect close connecion to the board
func (c *Client) Disconnect(ctx context.Context) (err error) {
	err = c.serialPort.Close()

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
	err = c.Connect(ctx)
	if err != nil {
		return err
	}

	c.Publish("reconnected", true)

	return nil
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(ctx context.Context, pin int, mode string) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debug("Pin: %d, Mode: %s", pin, mode)
		}

		if mode != client.ModeInput && mode != client.ModeInputPullup && mode != client.ModeOutput {
			return errors.Errorf("Can't found mode %s", mode)
		}

		url := fmt.Sprintf("/mode/%d/%s\n\r", pin, mode)

		com := c.read(ctx)
		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return ctx.Err()
		case err := <-com.Err:
			return err
		case resp = <-com.Res:

		}

		if c.isDebug {
			arest.Debug("Resp: %s", resp)
		}

		return nil
	}
}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(ctx context.Context, pin int, level int) (err error) {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if c.isDebug {
			log.Debugf("Pin: %d, Level: %s", pin, level)
		}

		if level != client.LevelHigh && level != client.LevelLow {
			return errors.Errorf("Can't found level %d", level)
		}

		url := fmt.Sprintf("/digital/%d/%d\n\r", pin, level)

		com := c.read(ctx)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return ctx.Err()
		case err := <-com.Err:
			return err
		case resp = <-com.Res:
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp)
		}

		return nil
	}
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(ctx context.Context, pin int) (level int, err error) {
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

		com := c.read(ctx)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return level, err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return level, ctx.Err()
		case err := <-com.Err:
			return level, err
		case resp = <-com.Res:
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

		com := c.read(ctx)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return nil, err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return nil, ctx.Err()
		case err := <-com.Err:
			return nil, err
		case resp = <-com.Res:
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
	select {
	case <-ctx.Done():
		return values, ctx.Err()
	default:
		c.mutex.Lock()
		defer c.mutex.Unlock()
		url := "/\n\r"
		data := make(map[string]interface{})

		com := c.read(ctx)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return nil, err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return nil, ctx.Err()
		case err := <-com.Err:
			return nil, err
		case resp = <-com.Res:
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

		com := c.read(ctx)

		_, err = c.serialPort.Write([]byte(url))
		if err != nil {
			return value, err
		}

		var resp string
		select {
		case <-ctx.Done():
			com.Cancel()
			return value, ctx.Err()
		case err := <-com.Err:
			return value, err
		case resp = <-com.Res:
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
