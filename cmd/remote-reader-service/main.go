package main

import (
	"context"
	"github.com/inconshreveable/log15"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/readers"
	"github.com/vshapovalov/rfid-reader-service/internal/services"
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const maxBuzzerInARow = 3

func main() {
	var err error

	appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	configService := services.NewFileConfigService(appDir)
	if err != nil {
		panic(err)
	}
	config, err := configService.GetConfig()
	if err != nil {
		panic(err)
	}

	logger := infrastructure.NewLog15Logger(config.IsDebugModeEnabled, false, []log15.Handler{
		log15.StreamHandler(&lumberjack.Logger{
			Filename:   "./logs/communication.log",
			MaxSize:    5,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   false,
		}, log15.LogfmtFormat()),
	})

	logger.Info("application started")

	var readerModule readers.IReaderModule

	for {
		readerModule, err = readers.CreateReader(config.Device, logger)
		if err != nil {
			logger.Error("cannot open device", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	logger.Info("device opened")

	var mqttBroker *infrastructure.MqttBroker

	for {
		mqttBroker, err = infrastructure.NewMqttBroker(
			logger,
			"reader-"+config.Id, config.MqttBroker.URI, config.MqttBroker.Username, config.MqttBroker.Password,
			true, utils.GetStatusTopic(config.Id),
			services.NewStatusMessage(config.Id, services.StatusOffline).ToByteArray(),
		)
		if err != nil {
			logger.Error("failed to create mqtt broker", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	logger.Info("broker connected")

	readerService := services.NewReaderService(
		readerModule, logger,
		config.ReverseCardNumber,
		config.UseBuzzerOnRead,
		maxBuzzerInARow,
		config.CardReadingInterval,
		config.ReadCards,
	)
	communicationService := services.NewBrokerCommunicationService(mqttBroker, config.Id, logger)
	logger.Info("services created")

	err = communicationService.Register(readerModule.GetReaderInfo())
	if err != nil {
		logger.Crit("failed to register service", "error", err)
		return
	}
	logger.Info("services registered")

	ctx, cancelFunc := context.WithCancel(context.Background())

	go readerService.ReadCards()

	err = communicationService.OnBuzzRequest(func(count int) {
		readerService.Buzz(count)
	})
	if err != nil {
		logger.Crit("failed to subscribe on buzzer message", "error", err)
		return
	}

	go func() {
		logger.Info("main loop started")
	mainLoop:
		for {
			select {
			case cardNumbers, ok := <-readerService.OnCardRead():
				logger.Info("service read card", "card", cardNumbers)
				if ok {
					err = communicationService.SendCardNumber(cardNumbers)
					if err != nil {
						logger.Error("failed to send card number", "error", err)
					}
					logger.Info("service sent card", "card", cardNumbers)
				} else {
					logger.Error("channel is closed")
				}
			case <-ctx.Done():
				logger.Info("main loop cancelling")
				break mainLoop
			}
		}
		logger.Info("main loop destroyed")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	osSignal := <-c
	logger.Info("application stopped", "signal", osSignal)
	cancelFunc()

	logger.Info("stopping reader service")
	<-readerService.StopCardsReading()
	logger.Info("reader service stopped")

	logger.Info("services unregister attempt")
	err = communicationService.Unregister()
	if err != nil {
		logger.Error("failed to unregister service", "error", err)
		return
	}
	logger.Info("services unregistered")

	logger.Info("closing mqtt broker")
	mqttBroker.Close()
	logger.Info("mqtt broker closed")

	logger.Info("closing reader module")
	err = readerModule.Close()
	if err != nil {
		logger.Error("cannot close reader", "error", err)
	}
	logger.Info("reader module closed")
}
