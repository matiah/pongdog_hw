package main

import (
	"fmt"
	io "godog/io_cardreader"
	mqtt "godog/mqtt"
)

func main() {

	transmitchan := make(chan mqtt.Mqtt_message, 10)

	io.CheckIfRoot()
	HIDPath := io.CheckForCardreader()
	cardreader := io.ConnectToCardReader(HIDPath)
	go io.ReadFromCardReaderAndTransmit(*cardreader, transmitchan)
	for {
		select {
		case x := <-transmitchan:
			fmt.Printf("Received from chan, topic: %s, message %s.\n", x.Topic, x.Message)
		}
	}
	/*fmt.Println(cardreader.DeviceInfo)
	fmt.Println("new-code")
	devicelist := hid.Enumerate(0, 0)
	fmt.Println(devicelist)
	kortleser, err := devicelist[0].Open()
	if err != nil {
		fmt.Println(err)
	}
	b := make([]byte, 3)
	for {
		_, err := kortleser.Read(b)
		if err != nil {
			fmt.Println(err)
		}
		if len(b) > 0 {
			fmt.Println(b)
		}
	}
	*/

	//var testLeser hid.Device

	//testLeser.DeviceInfo = io.EnumerateHID()
	//fmt.Println(testLeser)
	//testLeser.Open()

	//publishchan := make(chan mqtt.Mqtt_message)
	//go mqtt.StartMQTTserver(publishchan)

	//message := mqtt.Mqtt_message{
	//	Topic:   "input/card",
	//	Message: "1239591",
	//}

	//for i := 0; i < 10; i++ {
	//	publishchan <- message
	//	time.Sleep(time.Second)
	//}
}
