package restClient

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/gobot-arest/v1/plateforms/client"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
)

// Client implement arest interface
type Client struct {
	resty   *resty.Client
	isDebug bool
	url     string
	timeout time.Duration
	gobot.Eventer
}

// NewClient permit to initialize new client Object
func NewClient(url string, timeout time.Duration, isDebug bool) *Client {

	return &Client{
		resty:   nil,
		isDebug: isDebug,
		url:     url,
		timeout: timeout,
		Eventer: gobot.NewEventer(),
	}
}

// Client permit to get curent resty client
func (c *Client) Client() *resty.Client {
	return c.resty
}

// Connect start connection to the board
// We just try read port 0 value to check http connexion is ready
func (c *Client) Connect(ctx context.Context) (err error) {
	resty := resty.New().
		SetHostURL(c.url).
		SetHeader("Content-Type", "application/json")

	if c.timeout != 0 {
		resty.SetTimeout(c.timeout)
	}

	_, err = c.DigitalRead(ctx, 0)
	return err
}

// Disconnect close connecion to the board
func (c *Client) Disconnect(ctx context.Context) (err error) {
	c.resty = nil
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

	return nil
}

// SetPinMode permit to set pin mode
func (c *Client) SetPinMode(ctx context.Context, pin int, mode string) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if c.isDebug {
			log.Debugf("Pin: %d, Mode: %d", pin, mode)
		}

		if mode != client.ModeInput && mode != client.ModeInputPullup && mode != client.ModeOutput {
			return errors.Errorf("Can't found mode %s", mode)
		}

		url := fmt.Sprintf("/mode/%d/%s", pin, mode)

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			Post(url)

		if c.isDebug {
			log.Debugf("Resp: %s", resp.String())
		}

		return err
	}
}

// DigitalWrite permit to set level on pin
func (c *Client) DigitalWrite(ctx context.Context, pin int, level int) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if c.isDebug {
			log.Debugf("Pin: %d, Level: %d", pin, level)
		}

		if level != client.LevelHigh && level != client.LevelLow {
			return errors.Errorf("Can't found level %d", level)
		}

		url := fmt.Sprintf("/digital/%d/%d", pin, level)

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			Post(url)

		if c.isDebug {
			log.Debugf("Resp: %s", resp.String())
		}

		return err
	}
}

// DigitalRead permit to read level from pin
func (c *Client) DigitalRead(ctx context.Context, pin int) (level int, err error) {

	select {
	case <-ctx.Done():
		return level, ctx.Err()
	default:

		if c.isDebug {
			log.Debugf("Pin: %d", pin)
		}

		url := fmt.Sprintf("/digital/%d", pin)
		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			SetResult(&data).
			Get(url)
		if err != nil {
			return level, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s, %+v", resp.String(), data)
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
		if c.isDebug {
			log.Debugf("Value name: %s", name)
		}

		url := fmt.Sprintf("/%s", name)
		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			SetResult(&data).
			Get(url)
		if err != nil {
			return nil, err
		}

		log.Debugf("Resp: %s", resp.String())

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

		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			SetResult(&data).
			Get("/")
		if err != nil {
			return nil, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp.String())
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

		if c.isDebug {
			log.Debugf("Function: %s, param: %s", name, param)
		}

		url := fmt.Sprintf("/%s", name)

		data := make(map[string]interface{})

		resp, err := c.resty.R().
			SetQueryParams(map[string]string{
				"params": param,
			}).
			SetHeader("Accept", "application/json").
			SetContext(ctx).
			SetResult(&data).
			Post(url)
		if err != nil {
			return value, err
		}

		if c.isDebug {
			log.Debugf("Resp: %s", resp.String())
		}

		if temp, ok := data["return_value"]; ok {
			value = int(temp.(float64))
		} else {
			errors.Errorf("Function %s not found", name)
		}

		return value, err
	}
}
