package middleware

import (
	"context"
	"hw7/storage"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func InitLoggingInterceptor(logStorage *storage.LogStorage, subStorage *storage.SubscribersStorage) (grpc.UnaryServerInterceptor, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			p = &peer.Peer{Addr: nil}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		consumer := "unknown"
		if ok {
			consumers := md.Get("consumer")
			if len(consumers) > 0 {
				consumer = consumers[0]
			}
		}

		resp, err := handler(ctx, req)
		time.Sleep(10 * time.Millisecond)

		msg := storage.LogMessage{
			Timestamp: time.Now().Unix(),
			Consumer:  consumer,
			Method:    info.FullMethod,
			Host:      p.Addr.String(),
		}
		subStorage.NotifyAll(msg)
		logStorage.NotifyAll(msg)

		return resp, err
	}, nil
}

func InitLoggingStreamInterceptor(logStorage *storage.LogStorage, subStorage *storage.SubscribersStorage) (grpc.StreamServerInterceptor, error) {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		p, ok := peer.FromContext(ctx)
		host := "unknown"
		if ok && p != nil && p.Addr != nil {
			host = p.Addr.String()
		}

		md, ok := metadata.FromIncomingContext(ctx)
		consumer := "unknown"
		if ok {
			consumers := md.Get("consumer")
			if len(consumers) > 0 {
				consumer = consumers[0]
			}
		}

		msg := storage.LogMessage{
			Timestamp: time.Now().Unix(),
			Consumer:  consumer,
			Method:    info.FullMethod,
			Host:      host,
		}
		subStorage.NotifyAll(msg)
		logStorage.NotifyAll(msg)

		return handler(srv, ss)
	}, nil
}
