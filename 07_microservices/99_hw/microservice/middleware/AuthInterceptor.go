package middleware

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"

	"google.golang.org/grpc/metadata"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InitAuthInterceptor(acl string) (grpc.UnaryServerInterceptor, error) {
	var perms map[string][]string
	if err := json.Unmarshal([]byte(acl), &perms); err != nil {
		return nil, err
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "No metadata found")
		}
		vals := md.Get("consumer")
		if len(vals) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Consumer not set")
		}
		consumer := vals[0]

		consumerPermissions, exists := perms[consumer]
		if !exists {
			return nil, status.Error(codes.Unauthenticated, "Consumer not registered")
		}

		for _, perm := range consumerPermissions {
			match, err := filepath.Match(perm, info.FullMethod)
			if err != nil {
				log.Println("Error matching permission:", err)
				continue
			}
			if match {
				return handler(ctx, req)
			}
		}

		log.Println(info.FullMethod, "not allowed for consumer", consumer)
		return nil, status.Error(codes.Unauthenticated, "Consumer not allowed to use this method")
	}, nil
}

func InitAuthStreamInterceptor(acl string) (grpc.StreamServerInterceptor, error) {
	var perms map[string][]string
	if err := json.Unmarshal([]byte(acl), &perms); err != nil {
		return nil, err
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Unauthenticated, "No metadata found")
		}
		vals := md.Get("consumer")
		if len(vals) == 0 {
			return status.Error(codes.Unauthenticated, "Consumer not set")
		}
		consumer := vals[0]

		consumerPermissions, exists := perms[consumer]
		if !exists {
			return status.Error(codes.Unauthenticated, "Consumer not registered")
		}

		for _, perm := range consumerPermissions {
			match, err := filepath.Match(perm, info.FullMethod)
			if err != nil {
				log.Println("Error matching permission:", err)
				continue
			}
			if match {
				return handler(srv, ss)
			}
		}

		log.Println(info.FullMethod, "not allowed for consumer", consumer)
		return status.Error(codes.Unauthenticated, "Consumer not allowed to use this method")
	}, nil
}
