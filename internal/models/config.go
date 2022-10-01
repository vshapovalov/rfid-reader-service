package models

import "time"

type Config struct {
	Id                 string           `json:"id"`
	IsDebugModeEnabled bool             `json:"isDebugModeEnabled"`
	ReverseCardNumber  bool             `json:"reverseCardNumber"`
	UseBuzzerOnRead    bool             `json:"useBuzzerOnRead"`
	Device             ConfigDevice     `json:"device"`
	MqttBroker         ConfigMqttBroker `json:"mqttBroker"`
}

type ConfigDevice struct {
	Driver       string             `json:"driver"`
	M302Settings ConfigM302Settings `json:"M302Settings"`
}

type ConfigMqttBroker struct {
	URI      string `json:"uri"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ConfigM302Settings struct {
	Port        string        `json:"port"`
	Baud        int           `json:"baud"`
	ReadTimeout time.Duration `json:"readTimeout"`
	Size        byte          `json:"size"`
}
