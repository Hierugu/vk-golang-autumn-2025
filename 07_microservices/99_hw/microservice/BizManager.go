package main

import "context"

type BizManager struct {
	UnimplementedBizServer
}

func (bm BizManager) Add(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{Dummy: true}, nil
}

func (bm BizManager) Check(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{Dummy: true}, nil
}

func (bm BizManager) Test(ctx context.Context, req *Nothing) (*Nothing, error) {
	return &Nothing{Dummy: true}, nil
}
