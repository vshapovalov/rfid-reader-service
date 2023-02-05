package readers

import (
	"fmt"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/models"
	"github.com/vshapovalov/rfid-reader-service/internal/readers/drivers/rfidlib"
	"strings"
)

const driversPath = "drivers\\rfidlib\\drivers"

type RFIDLibReaderModule struct {
	logger infrastructure.ILogger
	reader *rfidlib.Reader
}

func NewRFIDLibReaderModule(deviceSettings models.ConfigRFIDLibSettings, logger infrastructure.ILogger) (*RFIDLibReaderModule, error) {
	drivers, err := rfidlib.LoadDrivers(driversPath)
	if err != nil {
		return nil, err
	}

	logger.Info("RFIDLib drivers loaded", "count", len(drivers))

	var readerDriver *rfidlib.Driver
	var targetDriver = strings.ToLower(deviceSettings.LibDriver)
	for _, driver := range drivers {
		if strings.ToLower(driver.Name) == targetDriver {
			readerDriver = driver
		}
	}

	if readerDriver == nil {
		return nil, fmt.Errorf("target lib driver not found: %s", deviceSettings.LibDriver)
	}

	var reader *rfidlib.Reader

	switch strings.ToLower(deviceSettings.Communication.Type) {
	case "com":
		comPort := deviceSettings.Communication.Settings["comPort"]
		if strings.TrimSpace(comPort) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication comPort: %s", comPort)
		}
		comBand := deviceSettings.Communication.Settings["comBand"]
		if strings.TrimSpace(comBand) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication comBand: %s", comBand)
		}
		frame := deviceSettings.Communication.Settings["frame"]
		if strings.TrimSpace(frame) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication frame: %s", frame)
		}
		reader = rfidlib.NewReader(readerDriver, rfidlib.ReaderCOMOptions{
			ComPort: comPort,
			ComBand: comBand,
			Frame:   frame,
		})
	case "usb":
		serialNumber := deviceSettings.Communication.Settings["serialNumber"]
		if strings.TrimSpace(serialNumber) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication serialNumber: %s", serialNumber)
		}

		hidItems, err := rfidlib.GetHIDItems(deviceSettings.LibDriver)
		if err != nil {
			return nil, err
		}

		var foundDevice bool
		for _, hidItem := range hidItems {
			if hidItem.SerialNum == serialNumber {
				foundDevice = true
				break
			}
		}

		if !foundDevice {
			return nil, fmt.Errorf("rfidlib reader device [%s] not found", serialNumber)
		}

		reader = rfidlib.NewReader(readerDriver, rfidlib.ReaderUSBOptions{
			AddrMode:     rfidlib.ReaderUSBAddrModeSerial,
			SerialNumber: serialNumber,
		})
	case "net":
		remoteIp := deviceSettings.Communication.Settings["remoteIp"]
		if strings.TrimSpace(remoteIp) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication remoteIp: %s", remoteIp)
		}
		remotePort := deviceSettings.Communication.Settings["remotePort"]
		if strings.TrimSpace(remotePort) == "" {
			return nil, fmt.Errorf("wrong rfidlib communication remotePort: %s", remotePort)
		}
		reader = rfidlib.NewReader(readerDriver, rfidlib.ReaderNETOptions{
			RemoteIp:   remoteIp,
			RemotePort: remotePort,
		})
	default:
		return nil, fmt.Errorf("unknow rfidlib communication type: %s", deviceSettings.Communication.Type)
	}

	err = reader.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open rfidlib reader: %w", err)
	}

	return &RFIDLibReaderModule{
		logger: logger,
		reader: reader,
	}, nil
}

func (m *RFIDLibReaderModule) Close() error {
	return m.reader.Close()
}

func (m *RFIDLibReaderModule) Buzz() error {
	return nil
}

func (m *RFIDLibReaderModule) ReadCards() ([][]byte, error) {
	return m.reader.ReadCards()
}
