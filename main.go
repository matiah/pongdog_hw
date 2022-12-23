package main

import (
	io_button "godog/io_buttons"
	io_card "godog/io_cardreader"
	mqtt "godog/mqtt"
)

func main() {
	transmitchan := make(chan mqtt.Mqtt_message, 10)
	// ---- system health checks ----
	io_card.CheckIfRoot()
	HIDPath := io_card.CheckForCardreader()
	cardreader := io_card.ConnectToCardReader(HIDPath)
	gpiochip := io_button.CheckForGPIOchip()
	// ---- start up goroutines ----
	go mqtt.StartMQTTserver(transmitchan)
	go io_card.ReadFromCardReaderAndTransmit(*cardreader, transmitchan)
	go io_button.PollButtonsAndTransmit(gpiochip, transmitchan)
	for {
		select {}
	}
}
