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

func (m *MetricStorage) Save(savingSettings *SavingSettings) {

	//if len(savingSettings.Database) > 0 {
	//	m.saveToDB(savingSettings.Database)
	//} else {
	m.saveToFile(savingSettings)
	//}

}
