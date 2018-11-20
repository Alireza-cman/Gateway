package client2gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"strings"

	"github.com/Atrovan/Gateway/flatten"
	"github.com/Atrovan/Gateway/variable"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/tidwall/sjson"
)

//Global variable
//TODO: I should by some how change this part of code
var PublishClient MQTT.Client

func init() {
	log.Println("trying to connect to the platform ")
	opts := MQTT.NewClientOptions().AddBroker(variable.PlatformBroker)
	PublishClient = MQTT.NewClient(opts)
	//log.Println("asdasd")

	if token := (PublishClient).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("[DONE] The gatewat is connected to the platform")

}

func Message2Platfrom(client MQTT.Client, value string, topic string, ShouldFlat bool) {

	result := value
	var err error
	if ShouldFlat == true {
		result, err = flatten.FlattenString(value, "", flatten.DotStyle)
	}
	//result, err := flatten.FlattenString(payloadMessage, "", flatten.DotStyle)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("output is: ", result)
	if token := client.Publish(topic, 0, false, result); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}

//D2G_connect: this function is dealing with the device to gateway connection
//v1/sensors/connect
var D2G_connect = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
		log.Println("something went wrong")
	}
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]
	if serialNumber == nil {
		log.Println("serialNumber is not found in the message ")
		return
	}
	serialNumber = fmt.Sprintf("%v", serialNumber)
	value, _ := sjson.Set(string(msg.Payload()[:]), "device", serialNumber)
	go Message2Platfrom(PublishClient, value, variable.G2P_Connect, true)

}

//v1/sensors/+/connect
var D2G_connect_v2 = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}

	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
	}

	RetrivedTopic := strings.Split(msg.Topic(), "/")
	fmt.Println(input)
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]
	//value := string(msg.Payload()[:])
	result, _ := sjson.Set("", "device", RetrivedTopic[2])
	//println(value)

	if serialNumber == nil {
		log.Println("serialNumber is not found in the message ")
		//value, _ = sjson.Set(string(msg.Payload()[:]), "serialNumber", RetrivedTopic[2])
	}
	//fmt.Println(value)
	go Message2Platfrom(PublishClient, result, variable.G2P_Connect, true)
}

//v1/sensors/disconnect
var D2G_disconnect = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
	}

	//RetrivedTopic := strings.Split(msg.Topic(), "/")
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]

	if serialNumber == nil {
		log.Println("serialNumber is not found in the message ")
		return
		//value, _ = sjson.Set(string(msg.Payload()[:]), "serialNumber", RetrivedTopic[2])
	}
	serialNumber = fmt.Sprintf("%v", serialNumber)
	result, _ := sjson.Set("", "device", serialNumber)
	//fmt.Println(value)
	go Message2Platfrom(PublishClient, result, variable.G2P_Disconnect, true)
}

var D2G_disconnect_v2 = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}

	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
	}

	RetrivedTopic := strings.Split(msg.Topic(), "/")

	result, _ := sjson.Set("", "device", RetrivedTopic[2])

	go Message2Platfrom(PublishClient, result, variable.G2P_Disconnect, true)
}

var D2G_sensors = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {

		log.Println(err)
		return
	}
	//telemetry := fmt.Sprintf("%v", string(msg.Payload()))

	output := input.(map[string]interface{})
	tmp := output["serialNumber"]
	if tmp == nil {
		log.Println(" serial number does not exist")
		return
	}
	//serialNumber := fmt.Sprintf("%v", tmp)
	serialNumber := tmp.(string)
	//fmt.Println(serialNumber)

	serialNumber = ":" + serialNumber
	total, _ := sjson.Set("", serialNumber+".0.ts", time.Now().UTC().Unix())

	fmt.Println(total)
	for k, v := range output {
		if k == "serialNumber" {
			continue
		}
		total, _ = sjson.Set(total, serialNumber+".0.value."+k, v)
	}
	fmt.Println(total)
	go Message2Platfrom(PublishClient, total, variable.G2P_Disconnect, false)
}

var D2G_sensors_v2 = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
		return
	}

	RetrivedTopic := strings.Split(msg.Topic(), "/")
	serialNumber := RetrivedTopic[2]

	//
	serialNumber = ":" + serialNumber
	//
	telemetryKey := RetrivedTopic[3]
	output := input.(map[string]interface{})
	//json_output, err := json.Marshal(output)
	value := output["value"]
	total, _ := sjson.Set("", serialNumber+".0.ts", time.Now().UTC().Unix())
	result, _ := sjson.Set(total, serialNumber+".0.value."+telemetryKey, value)
	fmt.Println(total)

	go Message2Platfrom(PublishClient, result, variable.G2P_Telemetry, false)
}

var D2G_RPC = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Println(err)
		return
	}

	RetrivedTopic := strings.Split(msg.Topic(), "/")
	serialNumber := RetrivedTopic[2]

	//
	serialNumber = ":" + serialNumber

	//
	telemetryKey := RetrivedTopic[3]
	output := input.(map[string]interface{})
	//json_output, err := json.Marshal(output)
	value := output["value"]
	total, _ := sjson.Set("", serialNumber+".0.ts", time.Now().UTC().Unix())
	result, _ := sjson.Set(total, serialNumber+".0.value."+telemetryKey, value)
	fmt.Println(total)

	go Message2Platfrom(PublishClient, result, variable.G2P_Telemetry, false)
}
