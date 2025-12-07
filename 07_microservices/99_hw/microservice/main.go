package main

func main() {
	println("usage: go test -v -race")

	// ctx, cancel := context.WithCancel(context.Background())
	// go StartMyMicroservice(ctx, ":8081", acl)

	// grpcClient, err := grpc.NewClient(":8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("cant connect to grpc: %v", err)
	// }

	// ad := NewBizClient(grpcClient)
	// nothing, err := ad.Add(context.Background(), &Nothing{})
	// log.Println("gRPC Resp:", nothing, err)

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	cancel()
	// 	time.Sleep(1 * time.Second)
	// 	wg.Done()
	// }()
	// wg.Wait()
}
