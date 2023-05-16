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
	logger = log.New(&buf, "error: ", log.Lshortfile)
	return logger
}
