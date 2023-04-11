package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func (m *MetricStorage) saveToFile(fileName string) {

	producer, err := newProducer(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		err = producer.close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	if err := producer.writeEvent(m); err != nil {
		fmt.Println(err.Error())
		return
	}
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

func (p *producer) writeEvent(metStorage *MetricStorage) error {
	return p.encoder.Encode(&metStorage)
}

func (p *producer) close() error {
	return p.file.Close()
}
