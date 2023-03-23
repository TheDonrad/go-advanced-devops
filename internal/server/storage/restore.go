package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

func (m *MetricStorage) Restore(restore bool, savingSettings *SavingSettings, metStorage *MetricStorage) {
	if !restore {
		return
	}
	fileName := savingSettings.StoreFile

	consumer, err := newConsumer(fileName)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		err = consumer.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	err = consumer.readEvent(metStorage)
	if err != nil {
		fmt.Println(err.Error())
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
