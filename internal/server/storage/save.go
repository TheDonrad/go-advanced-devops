package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type SavingSettings struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
}

func NewSavingSettings(storeInterval time.Duration, storeFile string) *SavingSettings {
	return &SavingSettings{
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
	}
}

func (m *MetricStorage) Save(savingSettings *SavingSettings) {

	fileName := savingSettings.StoreFile
	Producer, err := NewProducer(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		err = Producer.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	if err := Producer.WriteEvent(m); err != nil {
		fmt.Println(err.Error())
		return
	}
}

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(fileName string) (*Producer, error) {
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(fileName, flag, 0777)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(metStorage *MetricStorage) error {
	return p.encoder.Encode(&metStorage)
}

func (p *Producer) Close() error {
	return p.file.Close()
}
