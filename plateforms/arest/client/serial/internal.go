package serialClient

import (
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Com permit communication with read routine
type Com struct {
	Err      chan error
	Res      chan string
	Watchdog chan bool
}

func (c *Client) readProcess(ctx context.Context) {

	// Channel to sync watchdog with read routine
	chPing := make(chan bool)
	chEnd := make(chan bool)

	// Run watchdog routine when write
	go func() {
		for {
			// Start watchdog when write is called
			<-c.com.Watchdog

			timer := time.NewTicker(c.timeout)
			stopWatchdogLoop := false

			for !stopWatchdogLoop {
				select {
				case <-ctx.Done():
					// context done
					if c.isDebug {
						log.Debugf("Watchdog exit because of context done")
					}
					return
				case <-chPing:
					// Reset timer
					timer = time.NewTicker(c.timeout)
				case <-timer.C:
					// Timeout
					if c.isDebug {
						log.Debug("Watchdog detect timeout, we close connexion")
					}
					c.Publish("timeout", true)
					return
				case <-chEnd:
					// read complet
					stopWatchdogLoop = true
				}

				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	// Read routine
	go func() {
		buffer := make([]byte, 4096)
		var resp strings.Builder
		loop := true

		for {
			for loop {
				select {
				case <-ctx.Done():
					return
				default:
					n, err := c.serialPort.Read(buffer)
					if err != nil {
						c.com.Err <- err
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

					chPing <- true
				}
			}

			chEnd <- true
			loop = true
			c.com.Res <- resp.String()

			resp.Reset()
		}
	}()
}

// write permit to sync the read/write on serial
func (c *Client) write(ctx context.Context, url string) (res string, err error) {

	// Start watchdog
	c.com.Watchdog <- true

	// Write query on serial
	_, err = c.serialPort.Write([]byte(url))
	if err != nil {
		return "", err
	}

	// Wait result
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case err := <-c.com.Err:
		return "", err
	case res = <-c.com.Res:
	}

	return res, nil
}
