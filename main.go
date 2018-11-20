package main

import (
	"log"
	"sync"

	"github.com/Atrovan/Gateway/variable"

	"github.com/Atrovan/Gateway/client2gateway"
	"github.com/Atrovan/Gateway/platform2gateway"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func GW_Publish(payload []byte, client MQTT.Client) {

	if token := client.Publish("alireza/kesha", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func main() {
	log.Println("gateway is running")
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// device to gateway connection
	if token := client.Subscribe(variable.D2G_Connect, 0, client2gateway.D2G_connect); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Subscribe(variable.D2G_Connect_embed, 0, client2gateway.D2G_connect_v2); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// device to gatewat disconnection
	if token := client.Subscribe(variable.D2G_Disconnect, 0, client2gateway.D2G_disconnect); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Subscribe(variable.D2G_Disconnect_embed, 0, client2gateway.D2G_disconnect_v2); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// sensor
	if token := client.Subscribe(variable.D2G_Sensor, 0, client2gateway.D2G_sensors); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(variable.D2G_Sensor_embed, 0, client2gateway.D2G_sensors_v2); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// this portion make the connection and subscrtption alive
	if token := client.Subscribe(variable.G2P_RPC, 0, platform2gateway.G2P_RPC); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}
