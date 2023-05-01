package logs

import (
	"bytes"
	"log"
)

func New() *log.Logger {
	var buf bytes.Buffer
	logger := log.New(&buf, "error: ", log.Lshortfile)
	return logger
}
