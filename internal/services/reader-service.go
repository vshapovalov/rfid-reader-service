package services

import (
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/readers"
	"strconv"
	"strings"
	"time"
)

const cardReadingInterval = 100 * time.Millisecond
const actionReadCard = "readCard"
const actionBuzz = "buzz"

type ReaderService struct {
	CardReadInfo      chan string
	actions           chan string
	logger            infrastructure.ILogger
	readerModule      readers.IReaderModule
	reverseCardNumber bool
	useBuzzer         bool
	maxBuzzerInARow   int
	lastReadCard      int
}

func (r *ReaderService) Buzz(count int) {
	r.actions <- actionBuzz + ":" + strconv.Itoa(count)
}

func (r *ReaderService) readCard() {
	var targetCard int
	card, err := r.readerModule.ReadCard()
	if err != nil {
		r.logger.Error("cannot read card", "error", err)
		r.lastReadCard = 0
		return
	}

	if card == nil && len(card) == 0 {
		r.lastReadCard = 0
		return
	}

	r.logger.Info("card read", "card", card)

	if r.reverseCardNumber {
		targetCard = byteArrayToInt(reverseArray(card))
	} else {
		targetCard = byteArrayToInt(card)
	}

	if targetCard == r.lastReadCard {
		return
	}

	r.lastReadCard = targetCard
	if r.useBuzzer {
		r.doBuzz(1)
	}
	go func(val int) { r.CardReadInfo <- strconv.Itoa(val) }(targetCard)
}

func (r *ReaderService) doBuzz(count int) {
	for i := 0; i < count; i++ {
		err := r.readerModule.Buzz()
		if err != nil {
			r.logger.Error("cannot buzz", "error", err)
			break
		}
	}
	r.logger.Info("buzzer used", "count", count)
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
			}
		}
	}
}

func (r *ReaderService) ReadCards() {
	r.logger.Info("service readCard loop started")
	for {
		time.Sleep(cardReadingInterval)
		r.actions <- actionReadCard
	}
}

func (r *ReaderService) OnCardRead() chan string {
	return r.CardReadInfo
}

func NewReaderService(readerModule readers.IReaderModule, logger infrastructure.ILogger, reverseCardNumber, useBuzzer bool, maxBuzzerInARow int) *ReaderService {

	service := &ReaderService{
		CardReadInfo: make(chan string, 1),
		actions:      make(chan string, 1),
		logger:       logger,
		readerModule: readerModule, reverseCardNumber: reverseCardNumber, useBuzzer: useBuzzer, maxBuzzerInARow: maxBuzzerInARow}

	go service.doActions()

	return service
}
