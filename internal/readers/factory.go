package readers

import (
	"errors"
	"github.com/vshapovalov/rfid-reader-service/internal/infrastructure"
	"github.com/vshapovalov/rfid-reader-service/internal/models"
)

const m302 = "M302"

func CreateReader(device models.ConfigDevice, logger infrastructure.ILogger) (IReaderModule, error) {
	switch device.Driver {
	case m302:
		module, err := NewM302ReaderModule(device.M302Settings, logger)
		return module, err
	}
	return nil, errors.New("unknown device driver")
}
