package services

import (
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/readers"
	"strconv"
	"time"
)

const cardReadingInterval = 100 * time.Millisecond

type ReaderService struct {
	CardReadInfo      chan string
	logger            infrastructure.ILogger
	readerModule      readers.IReaderModule
	reverseCardNumber bool
	useBuzzer         bool
}

func (r *ReaderService) readCards() {
	r.logger.Info("service readCard loop started")
	var targetCard int
	var lastCard int
	for {
		time.Sleep(cardReadingInterval)
		card, err := r.readerModule.ReadCard()
		if err != nil {
			r.logger.Error("cannot read card", "error", err)
			lastCard = 0
			continue
		}

		if card == nil && len(card) == 0 {
			lastCard = 0
			continue
		}

		r.logger.Info("card read", "card", card)

		if r.reverseCardNumber {
			targetCard = byteArrayToInt(reverseArray(card))
		} else {
			targetCard = byteArrayToInt(card)
		}

		if targetCard == lastCard {
			continue
		}

		lastCard = targetCard
		if r.useBuzzer {
			err = r.readerModule.Buzz()
			if err != nil {
				r.logger.Error("cannot buzz", "error", err)
			}
			r.logger.Info("buzzer used")
		}
		go func(val int) { r.CardReadInfo <- strconv.Itoa(val) }(targetCard)
	}
}

func (r *ReaderService) OnCardRead() chan string {
	return r.CardReadInfo
}

func NewReaderService(readerModule readers.IReaderModule, logger infrastructure.ILogger, reverseCardNumber, useBuzzer bool) *ReaderService {

	service := &ReaderService{
		CardReadInfo: make(chan string, 1),
		logger:       logger,
		readerModule: readerModule, reverseCardNumber: reverseCardNumber, useBuzzer: useBuzzer}

	go service.readCards()

	return service
}
