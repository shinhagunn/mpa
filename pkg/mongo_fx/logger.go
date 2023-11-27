package mongo_fx

import (
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBlogger struct {
	log *logrus.Logger
	mu  sync.Mutex
}

func (logger *DBlogger) Info(level int, msg string, keyandvalues ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	if options.LogLevel(level+1) == options.LogLevelDebug {
		logger.log.Debugf("message: %s value:%v", msg, keyandvalues)
	} else {
		logger.log.Infof("message: %s value:%v", msg, keyandvalues)
	}
}
func (logger *DBlogger) Error(err error, msg string, _ ...interface{}) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.log.Errorf("error: %v, message: %s\n", err, msg)
}
