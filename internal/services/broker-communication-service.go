package services

import (
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
)

type BrokerCommunicationService struct {
	broker        infrastructure.IBroker
	readerId      string
	registerTopic string
	cardReadTopic string
	buzzerTopic   string
}

func (b BrokerCommunicationService) SendCardNumber(cardNumber string) error {
	return b.broker.SendMessage(b.registerTopic, NewCardReadInfo(cardNumber, b.readerId).ToByteArray(), false)
}

func (b BrokerCommunicationService) Register() error {
	return b.broker.SendMessage(b.registerTopic, NewStatusInfo(b.readerId, StatusOnline).ToByteArray(), true)
}
func (b BrokerCommunicationService) Unregister() error {
	return b.broker.SendMessage(b.registerTopic, NewStatusInfo(b.readerId, StatusOffline).ToByteArray(), true)
}

func NewBrokerCommunicationService(broker infrastructure.IBroker, readerId string) *BrokerCommunicationService {
	return &BrokerCommunicationService{
		broker: broker, readerId: readerId,
		registerTopic: utils.GetStatusTopic(readerId),
		cardReadTopic: utils.GetCardReadTopic(readerId),
		buzzerTopic:   utils.GetBuzzerTopic(readerId),
	}
}
