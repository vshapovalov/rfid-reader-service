package services

import "github.com/vshapovalov/rfid-reader-service/internal/models"

type (
	IConfigService interface {
		GetConfig() (*models.Config, error)
	}

	ICommunicationService interface {
		SendCardNumber(cardNumber string) error
		Register() error
		Unregister() error
		OnBuzzRequest(handler func(count int)) error
	}

	IReaderService interface {
		OnCardRead() chan string
		Buzz(count int)
		ReadCards()
	}
)
