package main

import (
	"os"
	"sync"

	"github.com/Atrovan/Gateway/client2gateway"
	"github.com/Atrovan/Gateway/platform2gateway"
	"github.com/Atrovan/Gateway/variable"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	logging "github.com/op/go-logging"
)

func GW_Publish(payload []byte, client MQTT.Client) {

	if token := client.Publish("alireza/kesha", 0, false, payload); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {

	//
	//logger format, dont touch it
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
	//
	//go Gateway2Platform()

	go Things2Gateway()
	//

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}

func Things2Gateway() {
	// MQTT client which is dealing with the platform
	opts := MQTT.NewClientOptions().AddBroker(variable.GatewayBroker)
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	//

	go ThingsManipulator(client)

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}

//ThingsManipulator is a kind of load balancer of topics
func ThingsManipulator(client MQTT.Client) {
	log.Notice("Thingsmanipulator is started")
	if token := client.Subscribe(variable.D2G_Connect, 0, client2gateway.D2G_connect); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	if token := client.Subscribe(variable.D2G_Connect_embed, 0, client2gateway.D2G_connect_v2); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	// device to gatewat disconnection
	if token := client.Subscribe(variable.D2G_Disconnect, 0, client2gateway.D2G_disconnect); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	if token := client.Subscribe(variable.D2G_Disconnect_embed, 0, client2gateway.D2G_disconnect_v2); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}

	// sensor
	if token := client.Subscribe(variable.D2G_Sensor, 0, client2gateway.D2G_sensors); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}

	if token := client.Subscribe(variable.D2G_Sensor_embed, 0, client2gateway.D2G_sensors_v2); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	// this portion make the connection and subscrtption alive
	if token := client.Subscribe(variable.G2P_RPC, 0, platform2gateway.G2P_RPC); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}
