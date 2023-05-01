package filestorage

import (
	"encoding/json"
	"log"
	"os"

	"goAdvancedTpl/internal/fabric/metrics"
)

func (m *FileStorage) Restore() {

	consumer, err := newConsumer(m.Settings.StoreFile)

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
	err = consumer.readEvent(&m.Metrics)
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

func (c *consumer) readEvent(metStorage *metrics.MetricStorage) error {
	if err := c.decoder.Decode(metStorage); err != nil {
		return err
	}
	return nil
}
