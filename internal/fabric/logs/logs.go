package logs

import (
	"bytes"
	"log"
	"sync"
)

var (
	logger *log.Logger
	buf    bytes.Buffer
	once   sync.Once
)

// New возвращает объект для логирования
func New() *log.Logger {
	once.Do(
		func() {
			logger = log.New(&buf, "error: ", log.Lshortfile)
		})
	//buf.Reset()
	return logger
}
