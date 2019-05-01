package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/Atrovan/Gateway/client2gateway"
	"github.com/Atrovan/Gateway/variable"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	logging "github.com/op/go-logging"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var Mut sync.Mutex

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
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	//config.GetLocalClient()
	//go Gateway2Platform()
	go Things2Gateway()
	//

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}

func Things2Gateway() {
	// reading config file
	Mut.Lock()
	configFile := variable.ConfigFile
	plan, err := ioutil.ReadFile(configFile)
	configFile = string(plan)
	if err != nil {
		log.Error(err)
		return
	}
	Mut.Unlock()
	GatewayBroker := gjson.Get(configFile, "GatewayAddress").Str
	// MQTT client which is dealing with the platform
	opts := MQTT.NewClientOptions().AddBroker(GatewayBroker)
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
	//this portion make the connection and subscrtption alive
	if token := client.Subscribe(variable.G2P_RPC, 0, platform2gateway.G2P_RPC); token.Wait() && token.Error() != nil {
		log.Error(token.Error())
	}
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}
