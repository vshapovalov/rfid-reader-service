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
	}

	IReaderService interface {
		OnCardRead() chan string
	}
)
