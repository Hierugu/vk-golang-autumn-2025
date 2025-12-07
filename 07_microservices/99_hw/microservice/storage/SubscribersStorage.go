package storage

import "sync"

type SubscribersStorage struct {
	subscribers sync.Map
}

type LogMessage struct {
	Timestamp int64
	Consumer  string
	Method    string
	Host      string
}

func NewSubscribersStorage() *SubscribersStorage {
	return &SubscribersStorage{}
}

func (ss *SubscribersStorage) Subscribe() chan LogMessage {
	subscriberChan := make(chan LogMessage, 10)
	ss.subscribers.Store(subscriberChan, struct{}{})
	return subscriberChan
}

func (ss *SubscribersStorage) Unsubscribe(subscriberChan chan LogMessage) {
	ss.subscribers.Delete(subscriberChan)
}

func (ss *SubscribersStorage) NotifyAll(msg LogMessage) {
	ss.subscribers.Range(func(key, value any) bool {
		subscriberChan := key.(chan LogMessage)
		subscriberChan <- msg
		return true
	})
}
