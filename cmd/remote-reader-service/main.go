package main

import (
	"context"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/readers"
	"github.com/vshapovalov/rfid-reader-service/internal/services"
	"github.com/vshapovalov/rfid-reader-service/internal/utils"
	"os"
	"os/signal"
	"path/filepath"
)

//const generalUserPassword = "qPAX890JT6c&h0rJ&aBNB#mDhgzrgOT7"
//const generalUserName = "general_user"

func main() {
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

	logger := infrastructure.NewLog15Logger(config.IsDebugModeEnabled, false)

	readerModule, err := readers.CreateReader(config.Device, logger)
	if err != nil {
		logger.Crit("cannot open device", "error", err)
		return
	}
	defer func() {
		err := readerModule.Close()
		if err != nil {
			logger.Error("cannot close device", "error", err)
			return
		}
		logger.Info("device closed")
	}()

	logger.Info("device opened")

	mqttBroker, err := infrastructure.NewMqttBroker(
		logger,
		"reader-"+config.Id, config.MqttBroker.URI, config.MqttBroker.Username, config.MqttBroker.Password,
		true, utils.GetStatusTopic(config.Id),
		services.NewStatusInfo(config.Id, services.StatusOffline).ToByteArray(),
	)
	if err != nil {
		logger.Crit("failed to create mqtt broker", "error", err)
		return
	}
	defer mqttBroker.Close()
	logger.Info("broker connected")

	readerService := services.NewReaderService(readerModule, logger, config.ReverseCardNumber, config.UseBuzzerOnRead)
	communicationService := services.NewBrokerCommunicationService(mqttBroker, config.Id)
	logger.Info("services created")

	err = communicationService.Register()
	if err != nil {
		logger.Crit("failed to register service", "error", err)
		return
	}
	logger.Info("services registered")

	defer func() {
		err := communicationService.Unregister()
		if err != nil {
			logger.Error("failed to unregister service", "error", err)
			return
		}
		logger.Info("services unregistered")
	}()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		logger.Info("main loop started")
	mainLoop:
		for {
			select {
			case cardNumber, ok := <-readerService.OnCardRead():
				logger.Info("service read card", "card", cardNumber)
				if ok {
					err = communicationService.SendCardNumber(cardNumber)
					if err != nil {
						logger.Error("failed to send card number", "error", err)
					}
					logger.Info("service sent card", "card", cardNumber)
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
}
