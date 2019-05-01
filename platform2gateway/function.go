package platform2gateway

import (
	"fmt"
	"io/ioutil"

	"github.com/Atrovan/Gateway/flatten"
	"github.com/Atrovan/Gateway/variable"
	logging "github.com/op/go-logging"
	"github.com/tidwall/gjson"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var PublishClient MQTT.Client
var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	// reading config file
	configFile := variable.ConfigFile
	plan, err := ioutil.ReadFile(configFile)
	configFile = string(plan)
	if err != nil {
		log.Error(err)
		return
	}
	GatewayBroker := gjson.Get(configFile, "GatewayAddress").Str
	log.Warning("trying to connect to the gateway broker ")
	opts := MQTT.NewClientOptions().AddBroker(GatewayBroker)
	PublishClient = MQTT.NewClient(opts)

	if token := (PublishClient).Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	log.Warning("[DONE] The client is connected to the local broker")
}

func Message2Platfrom(client MQTT.Client, value string, topic string, ShouldFlat bool) {

	result := value

	var err error
	if ShouldFlat == true {
		result, err = flatten.FlattenString(value, "", flatten.DotStyle)
	}
	//result, err := flatten.FlattenString(payloadMessage, "", flatten.DotStyle)
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println("output is: ", result)
	fmt.Println("topic  is:", topic)
	if token := client.Publish(topic, 0, false, result); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}
