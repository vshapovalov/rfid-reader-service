package services

import (
	"encoding/json"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
)

type BrokerCommunicationService struct {
	broker        infrastructure.IBroker
	readerId      string
	registerTopic string
	cardReadTopic string
	buzzerTopic   string
	logger        infrastructure.ILogger
}

func (b *BrokerCommunicationService) OnBuzzRequest(handler func(count int)) error {
	return b.broker.Subscribe(utils.GetBuzzerTopic(b.readerId), func(message []byte) {
		b.logger.Info("buzz message received")
		msg := new(BuzzerCount)
		err := json.Unmarshal(message, msg)
		if err != nil {
			b.logger.Error("cannot unmarshal buzzer payload", "error", err)
		}
		go handler(msg.Count)
	})
}

func (b *BrokerCommunicationService) SendCardNumber(cardNumbers []string) error {
	return b.broker.SendMessage(b.cardReadTopic, NewCardReadInfo(cardNumbers, b.readerId).ToByteArray(), false)
}

func (b *BrokerCommunicationService) Register(extraInfo string) error {
	return b.broker.SendMessage(b.registerTopic, NewStatusWithExtraInfoMessage(b.readerId, StatusOnline, extraInfo).ToByteArray(), true)
}
func (b *BrokerCommunicationService) Unregister() error {
	return b.broker.SendMessage(b.registerTopic, NewStatusMessage(b.readerId, StatusOffline).ToByteArray(), true)
}

func NewBrokerCommunicationService(broker infrastructure.IBroker, readerId string, logger infrastructure.ILogger) *BrokerCommunicationService {
	return &BrokerCommunicationService{
		broker: broker, readerId: readerId,
		registerTopic: utils.GetStatusTopic(readerId),
		cardReadTopic: utils.GetCardReadTopic(readerId),
		buzzerTopic:   utils.GetBuzzerTopic(readerId),
		logger:        logger,
	}
}
