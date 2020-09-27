package serialClient

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Com permit communication with read routine
type Com struct {
	Err    chan error
	cancel chan bool
	Res    chan string
}

// Cancel permit to stop current routine
func (c *Com) Cancel() {
	c.cancel <- true
}

func (c *Client) read(ctx context.Context) *Com {

	com := &Com{
		Err:    make(chan error),
		Res:    make(chan string),
		cancel: make(chan bool),
	}

	ping := make(chan bool)
	end := make(chan bool)

	// Run watchdog routine
	go func() {
		timer := time.NewTicker(c.timeout)
		select {
		case <-ctx.Done():
			if c.isDebug {
				log.Debugf("Watchdog exit because of context done")
			}
			return
		case <-ping:
			// Reset timer
			timer = time.NewTicker(c.timeout)
		case <-timer.C:
			if c.isDebug {
				log.Debug("Watchdog detect timeout, we reconnect on board")
			}
			err := c.Reconnect(ctx)
			if err != nil {
				com.Err <- err
				return
			}
			if c.isDebug {
				log.Debug("Watchdog successfully reconnect on board")
			}
		case <-end:
			return
		case <-com.Err:
			return
		case <-com.cancel:
			return
		}
	}()

	// Read process
	go func() {
		select {
		case <-ctx.Done():
			com.Err <- ctx.Err()
		default:

			buffer := make([]byte, 2048)
			var resp strings.Builder
			loop := true

			for loop {
				select {
				case <-com.Err:
					return
				case <-com.cancel:
					return
				default:
					n, err := c.serialPort.Read(buffer)
					if err != nil {
						com.Err <- err
						return
					}
					if n == 0 {
						loop = false
						break
					}
					resp.Write(buffer[:n])

					if strings.Contains(string(buffer[:n]), "\n") {
						loop = false
						break
					}

					ping <- true
				}
			}

			end <- true
			com.Res <- resp.String()
			return
		}
	}()

	return com
}
