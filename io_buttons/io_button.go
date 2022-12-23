package io_button

import (
	"fmt"
	mqtt "godog/mqtt"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

// Channel for internal messages.
var internalChan = make(chan mqtt.Mqtt_message, 10)

func CheckForGPIOchip() string {
	cc := gpiod.Chips()
	return cc[0]
}

// Button interrupt handler.
func handler(evt gpiod.LineEvent) {
	go notify(evt.Offset, evt.Timestamp)
}

func PollButtonsAndTransmit(gpiochip string, outputChan chan<- mqtt.Mqtt_message) {
	l, _ := gpiod.RequestLine(gpiochip, rpi.J8p37, gpiod.AsInput, gpiod.WithPullUp, gpiod.AsActiveLow, gpiod.WithEventHandler(handler), gpiod.WithRisingEdge, gpiod.WithDebounce(time.Millisecond*20))
	fmt.Println("✓ Found GPIO-chip: " + l.Chip())
	for {
		select {
		case internalMessage := <-internalChan:
			outputChan <- internalMessage
		}
	}
}
func notify(line int, timestamp time.Duration) {
	fmt.Printf("%d - Rising edge detected, button %d\n", int(timestamp.Seconds()), line)
	var mqttmessage mqtt.Mqtt_message
	mqttmessage.Topic = "input/button"
	mqttmessage.Message = fmt.Sprint(line)
	internalChan <- mqttmessage
}
