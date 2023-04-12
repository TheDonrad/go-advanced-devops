package storage

import (
	"time"
)

type SavingSettings struct {
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	Database      string
}

func NewSavingSettings(storeInterval time.Duration, storeFile string, database string) *SavingSettings {
	return &SavingSettings{
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
		Database:      database,
	}
}

func (m *MetricStorage) Save(database string, file string) error {

	if len(database) > 0 {
		m.saveToDB(database)
	} else {
		m.saveToFile(file)
	}
	return nil
}
