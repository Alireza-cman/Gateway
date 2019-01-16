package client2gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/Atrovan/Gateway/flatten"
	"github.com/Atrovan/Gateway/variable"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	logging "github.com/op/go-logging"
	"github.com/tidwall/sjson"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

//Global variable
//TODO: I should by some how change this part of code
var PublishClient MQTT.Client

//var taskQueue rmq.Queue

func init() {
	//logger format, dont touch it
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)
	//
	log.Warning("trying to connect to the platform ")
	opts := MQTT.NewClientOptions().AddBroker(variable.PlatformBroker)
	opts.Username = "tptuViJ6lrDwoMbnKcfQ"
	opts.ClientID = "AtrovanGatewayThings"
	opts.KeepAlive = 60
	PublishClient = MQTT.NewClient(opts)
	//log.Println("asdasd")

	if token := (PublishClient).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Info("[DONE] The gatewat is connected to the platform")
	go GatewayDescriptor(PublishClient)
	//connecting gateway to the redis queue
	// log.Println("[Zooring...] trying to make a connection between redis queue and the gateway")
	// connection := rmq.OpenConnection("my service", "tcp", "localhost:6379", 1)
	// taskQueue = connection.OpenQueue("tasks")
	// log.Println("[Done] Gateway is connected to the redis queue")

}

func Message2Platfrom(client MQTT.Client, value string, topic string, ShouldFlat bool) {

	// taskQueue.StartConsuming(10, time.Second)
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
	log.Critical("output is: ", result)
	log.Critical(topic)
	if token := client.Publish(topic, 0, false, result); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

}

//D2G_connect: this function is dealing with the device to gateway connection
//v1/sensors/connect
var D2G_connect = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
	// delivery := "task payload"
	// taskQueue.Publish(delivery)
	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Error(err)

	}
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]
	if serialNumber == nil {
		log.Error("serialNumber is not found in the message ")
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
		log.Error(err)
	}

	RetrivedTopic := strings.Split(msg.Topic(), "/")
	fmt.Println(input)
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]
	//value := string(msg.Payload()[:])
	result, _ := sjson.Set("", "device", RetrivedTopic[2])
	//println(value)

	if serialNumber == nil {
		log.Error("serialNumber is not found in the message ")
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
		log.Error(err)
	}

	//RetrivedTopic := strings.Split(msg.Topic(), "/")
	output := input.(map[string]interface{})
	serialNumber := output["serialNumber"]

	if serialNumber == nil {
		log.Error("serialNumber is not found in the message ")
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
		log.Error(err)
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

		log.Error(err)
		return
	}
	//telemetry := fmt.Sprintf("%v", string(msg.Payload()))

	output := input.(map[string]interface{})
	tmp := output["serialNumber"]
	if tmp == nil {
		log.Error(" serial number does not exist")
		return
	}
	//serialNumber := fmt.Sprintf("%v", tmp)
	serialNumber := tmp.(string)
	//fmt.Println(serialNumber)

	serialNumber = ":" + serialNumber
	total, _ := sjson.Set("", serialNumber+".0.ts", time.Now().UTC().Unix()*1000)

	fmt.Println(total)
	for k, v := range output {
		if k == "serialNumber" {
			continue
		}
		total, _ = sjson.Set(total, serialNumber+".0.values."+k, v)
	}
	fmt.Println(total)
	go Message2Platfrom(PublishClient, total, variable.G2P_Telemetry, false)
}

var D2G_sensors_v2 = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Error(err)
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
	result, _ := sjson.Set(total, serialNumber+".0.values."+telemetryKey, value)
	fmt.Println(total)

	go Message2Platfrom(PublishClient, result, variable.G2P_Telemetry, false)
}

var D2G_RPC = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	var input interface{}
	err := json.Unmarshal(msg.Payload(), &input)
	if err != nil {
		log.Error(err)
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

func GatewayDescriptor(client MQTT.Client) {
	aaaaa := `{"alireza":123123}`
	for {
		if token := client.Publish(variable.G2P_Dec, 0, false, aaaaa); token.Wait() && token.Error() != nil {

			log.Error(token.Error())
		} else {
			log.Notice("Gatway Descriptor was sent successfully to the Platform")
		}

		time.Sleep(time.Duration(60) * time.Second)
	}
}
