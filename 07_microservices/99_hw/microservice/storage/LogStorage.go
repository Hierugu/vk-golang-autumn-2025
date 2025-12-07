package storage

import "sync"

type LogStorage struct {
	loggers sync.Map
}

func NewLogStorage() *LogStorage {
	return &LogStorage{}
}

func (ls *LogStorage) Subscribe() chan LogMessage {
	loggerChan := make(chan LogMessage, 50)
	ls.loggers.Store(loggerChan, struct{}{})
	return loggerChan
}

func (ls *LogStorage) Unsubscribe(loggerChan chan LogMessage) {
	ls.loggers.Delete(loggerChan)
}

func (ls *LogStorage) NotifyAll(msg LogMessage) {
	ls.loggers.Range(func(key, value any) bool {
		loggerChan := key.(chan LogMessage)
		loggerChan <- msg
		return true
	})
}
