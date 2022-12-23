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

// Button interrupt handler.
func handler(evt gpiod.LineEvent) {
	go notify(evt.Offset, evt.Timestamp)
}

func PollButtonsAndTransmit(outputChan chan<- mqtt.Mqtt_message) {
	cc := gpiod.Chips()
	fmt.Println(cc)
	l, _ := gpiod.RequestLine("gpiochip0", rpi.J8p37, gpiod.AsInput, gpiod.WithPullUp, gpiod.AsActiveLow, gpiod.WithEventHandler(handler), gpiod.WithRisingEdge, gpiod.WithDebounce(time.Millisecond*20))
	fmt.Println(l)
	for {
		select {
		case internalMessage := <-internalChan:
			outputChan <- internalMessage
		}
	}
}
func notify(line int, timestamp time.Duration) {
	fmt.Printf("%s - Rising edge detected, button %d\n", timestamp, line)
	var mqttmessage mqtt.Mqtt_message
	mqttmessage.Topic = "input/button"
	mqttmessage.Message = fmt.Sprint(line)
	internalChan <- mqttmessage
}
