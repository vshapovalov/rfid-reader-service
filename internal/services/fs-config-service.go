package services

import (
	"encoding/json"
	"fmt"
	"github.com/vshapovalov/rfid-reader-service/internal/models"
	"io"
	"os"
)

const fileName = "config.json"

type FileConfigService struct {
	appDir string
}

func NewFileConfigService(appDir string) *FileConfigService {
	return &FileConfigService{appDir}
}

func (f FileConfigService) GetConfig() (*models.Config, error) {
	config := new(models.Config)

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", f.appDir, fileName), os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileData, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
