package storage

func (m *MetricStorage) Restore(restore bool, savingSettings *SavingSettings) {
	if !restore {
		return
	}
	if len(savingSettings.Database) == 0 {
		m.restoreFromFile(savingSettings.StoreFile)
	} else {
		m.restoreFromDB(savingSettings.Database)
	}

}
