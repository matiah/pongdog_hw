package io_cardreader

import (
	"fmt"
	mqtt "godog/mqtt"
	"log"
	"os"
	"strconv"

	hid "github.com/bearsh/hid"
)

// Checks whether the user is root / has SUDO privileges.
func CheckIfRoot() {
	euid := os.Geteuid()
	if euid != 0 {
		log.Fatal("Program must be run with root-privileges (sudo),euid = " + fmt.Sprint(euid) + ", wanted = 0. Exiting.")
	}
}

// Checks whether the RFID card reader is connected. Returns the path to the cardreader if found.
func CheckForCardreader() hid.DeviceInfo {
	devicelist := hid.Enumerate(65535, 53)
	if len(devicelist) == 0 {
		log.Fatal("Cardreader not detected. Exiting.")
	} else {
		fmt.Println("✓ Found cardreader: " + devicelist[0].Product)
	}
	return devicelist[0]
}

// Connects to the card reader and returns a handle that can be read from or closed.
func ConnectToCardReader(HIDPath hid.DeviceInfo) *hid.Device {
	cardreader, err := HIDPath.Open()
	if err != nil {
		log.Fatal("Error when opening cardreader: " + fmt.Sprint(err))
	}
	return cardreader

}

// Goroutine: Reads card data from the cardreader and sends it to the MQTT broker.
func ReadFromCardReaderAndTransmit(cardreader hid.Device, outputChan chan<- mqtt.Mqtt_message) {
	b := make([]byte, 3)
	var outputString string
	counter := 0
	for {
		size, err := cardreader.ReadTimeout(b, -1)
		if err != nil {
			fmt.Println(err)
			outputString = "0"
			counter = 0
		}
		if size > 0 && b[2] > 0 {
			if b[2] == 40 || counter == 10 {
				counter = 0
				fmt.Printf("-----------------\n Card detected: \n")
				fmt.Printf("Rawnumber = %s\n", outputString)
				output := reverseBytes(outputString)
				outputChan <- convertToMqttMessage(output)
				outputString = ""
			} else {
				if b[2] == 39 {
					outputString += "0"
				} else {
					outputString += strconv.Itoa(int(b[2] - 29))
				}
				counter++
			}
		}
	}
}

func reverseBytes(inputstring string) string {
	toInt, _ := strconv.Atoi(inputstring)
	toBinary := fmt.Sprintf("%032b", toInt)
	var output string
	for i := 0; i < 4; i++ {
		subslice := toBinary[i*8 : (i+1)*8]
		output += reverse(subslice)
	}
	result, _ := strconv.ParseInt(output, 2, 32)
	fmt.Printf("EM-number = %d\n", result)
	return strconv.Itoa(int(result))
}

func reverse(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}

func convertToMqttMessage(cardnumber string) mqtt.Mqtt_message {
	message := mqtt.Mqtt_message{
		Topic:   "input/card",
		Message: cardnumber,
	}
	return message
}
