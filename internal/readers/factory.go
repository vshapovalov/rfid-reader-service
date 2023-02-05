package readers

import (
	"errors"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/models"
	"strings"
)

const m302 = "M302"
const rfidLib = "RFIDLIB"

func CreateReader(device models.ConfigDevice, logger infrastructure.ILogger) (IReaderModule, error) {
	switch strings.ToUpper(device.Driver) {
	case m302:
		module, err := NewM302ReaderModule(device.M302Settings, logger)
		return module, err
	case rfidLib:
		module, err := NewRFIDLibReaderModule(device.RFIDLibSettings, logger)
		return module, err

	}
	return nil, errors.New("unknown device driver")
}
