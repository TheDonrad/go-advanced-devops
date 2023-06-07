package logs

import (
	"log"
	"os"
)

var (
	logger *log.Logger
)

// Logger возвращает объект для логирования
func Logger() *log.Logger {
	return logger
}

func init() {
	logger = log.New(os.Stdout, "error: ", log.Lshortfile)
}
