package tests

import (
	"context"
	"homework10/internal/ports"
	grpcPort "homework10/internal/ports/grpc"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var BenchSink int

func BenchmarkHTTP(b *testing.B) {
	ctx, cf := context.WithCancel(context.Background())
	endChan := make(chan int)
	hsrv, _ := ports.CreateServer(ctx, endChan)
	client := getTestClient(hsrv.Addr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.createUser("John", "john.doe@mail.com")
		if err != nil {
			panic("create user failed")
		}
		BenchSink++
	}
	cf()
	<-endChan
}

func BenchmarkGRPC(b *testing.B) {
	ctx, cf := context.WithCancel(context.Background())
	endChan := make(chan int)
	ports.CreateServer(ctx, endChan)
	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("connection failed")
	}
	client := grpcPort.NewAdServiceClient(conn)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
		if err != nil {
			panic("create user failed")
		}
		BenchSink++
	}
	cf()
	<-endChan
}
