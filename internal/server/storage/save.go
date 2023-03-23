package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type SavingSettings struct {
	StoreInterval int
	StoreFile     string
	Restore       bool
}

func NewSavingSettings(storeInterval int, storeFile string) *SavingSettings {
	return &SavingSettings{
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
	}
}

func (m *MetricStorage) Save(savingSettings *SavingSettings) {

	fileName := savingSettings.StoreFile
	producer, err := NewProducer(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer producer.Close()
	if err := producer.WriteEvent(m); err != nil {
		fmt.Println(err.Error())
		return
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*producer, error) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fileName, flag, 0777)
	if err != nil {
		return nil, err
	}
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteEvent(metStorage *MetricStorage) error {
	return p.encoder.Encode(&metStorage)
}

func (p *producer) Close() error {
	return p.file.Close()
}
