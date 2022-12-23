package mqtt_code

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Mqtt_message struct {
	Topic   string
	Message string
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}

// Establishes an MQTT-connection to Mosquitto (127.0.0.1:1883) and publishes any message sent to the publischan-channel. Messages must use the Mqtt_message struct format.
func StartMQTTserver(publishchan chan (Mqtt_message)) {
	var broker = "127.0.0.1:1883"
	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.SetClientID("goDog_mqtt")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case msg := <-publishchan:
			topic := msg.Topic
			message := msg.Message
			token = client.Publish(topic, 0, false, message)
			token.Wait()
		}
	}
}
