package storage

import (
	"encoding/json"
	"log"
	"os"
)

func (m *MetricStorage) restoreFromFile(fileName string) {

	consumer, err := newConsumer(fileName)

	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		err = consumer.Close()
		if err != nil {
			log.Print(err.Error())
		}
	}()
	err = consumer.readEvent(m)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func newConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (c *consumer) readEvent(metStorage *MetricStorage) error {
	if err := c.decoder.Decode(metStorage); err != nil {
		return err
	}
	return nil
}
