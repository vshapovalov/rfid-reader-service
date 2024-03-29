package services

import (
	"encoding/hex"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/readers"
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
	"strconv"
	"strings"
	"time"
)

const actionReadCard = "readCard"
const actionBuzz = "buzz"
const stopServiceLoop = "stopServiceLoop"

type ReaderService struct {
	CardReadInfo        chan []string
	actions             chan string
	serviceStopped      chan bool
	logger              infrastructure.ILogger
	readerModule        readers.IReaderModule
	reverseCardNumber   bool
	useBuzzer           bool
	maxBuzzerInARow     int
	readCards           bool
	lastReadCards       map[string]byte
	cardReadingInterval time.Duration
	isServiceStopped    bool
}

func (r *ReaderService) Buzz(count int) {
	r.actions <- actionBuzz + ":" + strconv.Itoa(count)
}

func (r *ReaderService) readCard() {
	r.logger.Debug("readCard")
	if r.isServiceStopped {
		return
	}
	r.logger.Debug("readerModule.ReadCards")
	cards, err := r.readerModule.ReadCards()
	if err != nil {
		r.logger.Error("cannot read card", "error", err)
		r.lastReadCards = make(map[string]byte)
		return
	}
	newCards := make([]string, 0)
	cardsMap := make(map[string]byte)
	r.logger.Debug("cards loop")
	for _, card := range cards {
		if len(card) == 0 {
			continue
		}
		var cardHex string
		if r.reverseCardNumber {
			cardHex = hex.EncodeToString(reverseArray(card))
		} else {
			cardHex = hex.EncodeToString(card)
		}

		cardsMap[cardHex] = 0

		if _, ok := r.lastReadCards[cardHex]; !ok {
			newCards = append(newCards, cardHex)
			r.logger.Info("card read", "card", cardHex)
		}
	}

	r.lastReadCards = cardsMap

	if len(newCards) > 0 {
		if r.useBuzzer {
			r.doBuzz(1)
		}

		go func(cards []string) { r.CardReadInfo <- cards }(newCards)
	}

}

func (r *ReaderService) doBuzz(count int) {
	if r.isServiceStopped {
		return
	}
	r.logger.Info("buzzer use started", "count", count)
	for i := 0; i < count; i++ {
		err := r.readerModule.Buzz()
		if err != nil {
			r.logger.Error("cannot buzz", "error", err)
			break
		}
	}
	r.logger.Info("buzzer use completed", "count", count)
}

func (r *ReaderService) doActions() {
	r.logger.Info("service actions loop started")
	for {
		select {
		case action := <-r.actions:
			actionParts := strings.Split(action, ":")
			switch actionParts[0] {
			case actionReadCard:
				r.readCard()
			case actionBuzz:
				count := 1
				if len(actionParts) == 2 {
					tmpCount, err := strconv.Atoi(actionParts[1])
					if err == nil && tmpCount > 0 && tmpCount <= r.maxBuzzerInARow {
						count = tmpCount
					}
				}
				r.doBuzz(count)
			case stopServiceLoop:
				r.serviceStopped <- true
			}
		}
	}
}

func (r *ReaderService) ReadCards() {
	if !r.readCards {
		r.logger.Warn("service readCard loop not started, reading cards disabled")
		return
	}
	r.logger.Info("service readCard loop started")
	for {
		if r.isServiceStopped {
			break
		}
		r.actions <- actionReadCard
		time.Sleep(r.cardReadingInterval)
	}
}

func (r *ReaderService) OnCardRead() chan []string {
	return r.CardReadInfo
}

func (r *ReaderService) StopCardsReading() chan bool {
	r.isServiceStopped = true
	r.actions <- stopServiceLoop
	return r.serviceStopped
}

func NewReaderService(readerModule readers.IReaderModule, logger infrastructure.ILogger, reverseCardNumber, useBuzzer bool, maxBuzzerInARow int, cardReadingInterval utils.Duration, readCards bool) *ReaderService {

	service := &ReaderService{
		CardReadInfo:        make(chan []string, 1),
		actions:             make(chan string, 1),
		logger:              logger,
		readerModule:        readerModule,
		reverseCardNumber:   reverseCardNumber,
		useBuzzer:           useBuzzer,
		maxBuzzerInARow:     maxBuzzerInARow,
		readCards:           readCards,
		cardReadingInterval: cardReadingInterval.Duration,
		lastReadCards:       make(map[string]byte),
		serviceStopped:      make(chan bool, 1),
	}

	logger.Info("Reader service starting", "readCards", readCards, "cardReadingInterval", cardReadingInterval, "reverseCardNumber", reverseCardNumber)

	go service.doActions()

	return service
}
