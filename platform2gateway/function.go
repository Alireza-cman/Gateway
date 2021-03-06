package platform2gateway

import (
	"encoding/json"
	"fmt"

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
	log.Warning("trying to connect to the gateway broker ")
	opts := MQTT.NewClientOptions().AddBroker(variable.GatewayBroker)
	PublishClient = MQTT.NewClient(opts)

	if token := (PublishClient).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
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

var G2P_RPC = func(client MQTT.Client, msg MQTT.Message) {
	//fmt.Printf("TOPIC: %s\n", msg.Topic())
	//fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {

		log.Error(err)
		return
	}
	//telemetry := fmt.Sprintf("%v", string(msg.Payload()))

	output := input.(map[string]interface{})
	tmp := output["device"]
	if tmp == nil {
		log.Error(" serial number does not exist")
		return
	}
	jsonInput, _ := json.Marshal(output)
	totalRPC := string(jsonInput)
	//fmt.Println(totalRPC)
	serialNumber := gjson.Get(totalRPC, "device")
	method := gjson.Get(totalRPC, "data.method")
	params := gjson.Get(totalRPC, "data.params")

	//	fmt.Println(params)
	topic := "v1/sensors/" + serialNumber.Str + "/request/" + method.Str + "/1"
	//fmt.Println(topic)

	go Message2Platfrom(PublishClient, params.Raw, topic, false)
}
