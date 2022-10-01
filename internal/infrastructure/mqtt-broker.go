package infrastructure

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MqttBroker struct {
	mqttClient mqtt.Client
	logger     ILogger
}

func (m *MqttBroker) Close() {
	m.mqttClient.Disconnect(250)
}

func (m *MqttBroker) SendMessage(topic string, message []byte, retained bool) error {
	token := m.mqttClient.Publish(topic, 0, retained, message)
	token.Wait()
	return token.Error()
}

func (m *MqttBroker) Subscribe(topic string, handler func(message []byte)) error {
	token := m.mqttClient.Subscribe(topic, 0, func(client mqtt.Client, message mqtt.Message) {
		handler(message.Payload())
	})
	token.Wait()
	return token.Error()
}

func NewMqttBroker(
	logger ILogger,
	clientId string, uri string, username string, password string,
	isWillEnabled bool, willTopic string, willPayload []byte,
) (*MqttBroker, error) {

	opts := mqtt.NewClientOptions().
		AddBroker(uri).
		SetClientID(clientId).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(1 * time.Second).
		SetUsername(username).
		SetPassword(password)

	if isWillEnabled {
		opts.WillEnabled = true
		opts.WillTopic = willTopic
		opts.WillPayload = willPayload
		opts.WillRetained = true
	}

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MqttBroker{mqttClient: c, logger: logger}, nil
}
