package models

import (
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
	"time"
)

type Config struct {
	Id                  string           `json:"id"`
	IsDebugModeEnabled  bool             `json:"isDebugModeEnabled"`
	ReverseCardNumber   bool             `json:"reverseCardNumber"`
	UseBuzzerOnRead     bool             `json:"useBuzzerOnRead"`
	CardReadingInterval utils.Duration   `json:"cardReadingInterval"`
	Device              ConfigDevice     `json:"device"`
	MqttBroker          ConfigMqttBroker `json:"mqttBroker"`
}

type ConfigDevice struct {
	Driver          string                `json:"driver"`
	M302Settings    ConfigM302Settings    `json:"M302Settings"`
	RFIDLibSettings ConfigRFIDLibSettings `json:"RFIDLibSettings"`
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

type ConfigRFIDLibSettings struct {
	LibDriver     string                       `json:"libDriver"`
	Communication RFIDLibSettingsCommunication `json:"communication"`
}

type RFIDLibSettingsCommunication struct {
	Type     string            `json:"type"`
	Settings map[string]string `json:"settings"`
}
