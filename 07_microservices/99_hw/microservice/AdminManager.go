package main

import (
	"fmt"
	"hw7/storage"
	"time"
)

type AdminManager struct {
	UnimplementedAdminServer
	logs        *storage.LogStorage
	subscribers *storage.SubscribersStorage
}

func (an AdminManager) Logging(n *Nothing, stream Admin_LoggingServer) error {
	logCh := an.subscribers.Subscribe()
	defer an.subscribers.Unsubscribe(logCh)

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case msg, ok := <-logCh:
			if !ok {
				return nil
			}
			event := &Event{
				Timestamp: msg.Timestamp,
				Consumer:  msg.Consumer,
				Method:    msg.Method,
				Host:      msg.Host,
			}
			if err := stream.Send(event); err != nil {
				return err
			}
		}
	}
}

func (an AdminManager) Statistics(n *StatInterval, stream Admin_StatisticsServer) error {
	ticker := time.NewTicker(time.Duration(n.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	logCh := an.logs.Subscribe()
	defer an.logs.Unsubscribe(logCh)

	byMethod, byConsumer := make(map[string]uint64), make(map[string]uint64)
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case logMsg, ok := <-logCh:
			if !ok {
				return nil
			}
			byMethod[logMsg.Method]++
			byConsumer[logMsg.Consumer]++
		case <-ticker.C:
			stat := &Stat{
				ByMethod:   byMethod,
				ByConsumer: byConsumer,
			}
			if err := stream.Send(stat); err != nil {
				return err
			}
			byMethod, byConsumer = make(map[string]uint64), make(map[string]uint64)
		}
	}
}
