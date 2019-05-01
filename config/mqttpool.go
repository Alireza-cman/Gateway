package config

import (
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var CloudClient *MqttClient
var LocalClient *MqttClient
var once sync.Once

type MqttClient struct {
	client      MQTT.Client
	subscribers map[string]MQTT.Token
}

func (r *MqttClient) AddSubscribtion(topic string, qos byte, function MQTT.MessageHandler) error {
	token := r.client.Subscribe(topic, qos, function)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	r.subscribers[topic] = token
	return nil
}
func (r *MqttClient) Publish(topic string, qos byte, message interface{}) error {
	token := r.client.Publish(topic, qos, false, message)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (r *MqttClient) RemoveSubscribtion(topic string) error {
	token := r.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func GetCloudClient(broker string) *MqttClient {
	var client MQTT.Client
	once.Do(func() {
		opts := MQTT.NewClientOptions().AddBroker(broker)
		client = MQTT.NewClient(opts)
		CloudClient = &MqttClient{client: client}
	})
	return CloudClient
}

func GetLocalClient(broker string) *MqttClient {
	var client MQTT.Client
	once.Do(func() {
		opts := MQTT.NewClientOptions().AddBroker(broker)
		client = MQTT.NewClient(opts)
		LocalClient = &MqttClient{client: client}
	})
	return LocalClient
}
