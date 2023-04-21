package filestorage

import (
	"encoding/json"
	"goAdvancedTpl/internal/fabric/metricsstorage"
	"log"
	"os"
)

func (m *FileStorage) Save() error {

	producer, err := newProducer(m.Settings.StoreFile)
	if err != nil {
		return err
	}

	defer func() {
		err = producer.close()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	m.Mutex.RLock()
	err = producer.writeEvent(&m.Metrics)
	m.Mutex.RUnlock()
	return err
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func newProducer(fileName string) (*producer, error) {
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

func (p *producer) writeEvent(metStorage *metricsstorage.MetricStorage) error {
	return p.encoder.Encode(&metStorage)
}

func (p *producer) close() error {
	return p.file.Close()
}
