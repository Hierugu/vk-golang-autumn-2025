package main

import (
	"context"
	"hw7/middleware"
	"hw7/storage"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc "google.golang.org/grpc"
)

func StartMyMicroservice(ctx context.Context, addr string, acl string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println("can't listen port", err)
		return err
	}
	logs := storage.NewLogStorage()
	subs := storage.NewSubscribersStorage()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	authInter, err := middleware.InitAuthInterceptor(acl)
	if err != nil {
		log.Println("can't init auth interceptor:", err)
		listener.Close()
		return err
	}

	authStreamInter, err := middleware.InitAuthStreamInterceptor(acl)
	if err != nil {
		log.Println("can't init auth stream interceptor:", err)
		listener.Close()
		return err
	}

	logInter, err := middleware.InitLoggingInterceptor(logs, subs)
	if err != nil {
		log.Println("can't init logging interceptor:", err)
		listener.Close()
		return err
	}

	logStreamInter, err := middleware.InitLoggingStreamInterceptor(logs, subs)
	if err != nil {
		log.Println("can't init logging stream interceptor:", err)
		listener.Close()
		return err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(authInter, logInter),
		grpc.ChainStreamInterceptor(authStreamInter, logStreamInter),
	)

	RegisterAdminServer(server, AdminManager{logs: logs, subscribers: subs})
	RegisterBizServer(server, BizManager{})

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- server.Serve(listener)
		close(serveErrCh)
	}()

	go func() {
		select {
		case <-serveErrCh:
			listener.Close()
			return
		case <-sigCh:
			server.GracefulStop()
		case <-ctx.Done():
			server.Stop()
		}

		if err := <-serveErrCh; err != nil {
			log.Printf("gRPC server exited with error after stop: %v\n", err)
		}
		listener.Close()
	}()

	return nil
}
