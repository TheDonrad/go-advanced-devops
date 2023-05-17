package logs

import (
	"log"
	"os"
	"sync"
)

var (
	logger *log.Logger
	once   *sync.Once
)

// Logger возвращает объект для логирования
func Logger() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "error: ", log.Lshortfile)
	})

	return logger
}
