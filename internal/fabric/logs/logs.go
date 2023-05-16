package logs

import (
	"bytes"
	"log"
)

var (
	logger *log.Logger
	buf    bytes.Buffer
)

// New возвращает объект для логирования
func New() *log.Logger {
	logger = log.New(&buf, "error: ", log.Lshortfile)
	return logger
}
