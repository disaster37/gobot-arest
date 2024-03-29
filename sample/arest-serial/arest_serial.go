package main

import (
	"context"
	"time"

	"github.com/disaster37/gobot-arest/plateforms/arest"
	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

func main() {

	log.SetLevel(log.DebugLevel)

	arestSerial := arest.NewSerialAdaptor("/dev/ttyUSB0", 5*time.Second, false)
	led := gpio.NewLedDriver(arestSerial, "3")

	// Input pullup button
	button := gpio.NewButtonDriver(arestSerial, "41")
	button.DefaultState = 1

	// Relay with normally closed
	relay := gpio.NewRelayDriver(arestSerial, "46")
	relay.Inverted = true

	// Put button as INPUT_PULLUP
	err := arestSerial.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = arestSerial.Board.SetPinMode(context.TODO(), 41, client.ModeInputPullup)
	if err != nil {
		log.Fatal(err)
	}

	work := func() {
		if err := led.Off(); err != nil {
			log.Error(err)
		}
		if err := relay.Off(); err != nil {
			log.Error(err)
		}
		if err := button.On(gpio.ButtonPush, func(s interface{}) {
			log.Debug("Pushed")
			err := led.Toggle()
			if err != nil {
				log.Error(err)
			}

			err = relay.Toggle()
			if err != nil {
				log.Error(err)
			}
		}); err != nil {
			panic(err)
		}
	}

	robot := gobot.NewRobot("arest",
		[]gobot.Connection{arestSerial},
		[]gobot.Device{led, button, relay},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}

}
