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

// New возвращает объект для логирования
func New() *log.Logger {
	once.Do(func() {
		logger = log.New(os.Stdout, "error: ", log.Lshortfile)
	})

	return logger
}
