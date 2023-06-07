package filestorage

import (
	"encoding/json"
	"log"
	"os"

	"goAdvancedTpl/internal/fabric/metrics"
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
	m.Mutex.Lock()
	err = producer.writeEvent(&m.Metrics)
	m.Mutex.Unlock()
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

func (p *producer) writeEvent(metStorage *metrics.MetricStorage) error {
	return p.encoder.Encode(&metStorage)
}

func (p *producer) close() error {
	return p.file.Close()
}
