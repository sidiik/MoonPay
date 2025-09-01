package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once // logger will be initialized once
)

func Init() {
	once.Do(func() {
		var err error
		log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	})
}

func L() *zap.Logger {
	if log == nil {
		Init()
	}

	return log
}
