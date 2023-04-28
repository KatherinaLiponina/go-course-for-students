package tests

import (
	"context"
	"log"
	"testing"

	"homework9/internal/ports"
	grpcPort "homework9/internal/ports/grpc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestGRPCCreateUser(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, "Oleg", res.Name)

	cancelFunc()
}

func TestGRPCGetUser(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	usr, err := client.GetUser(context.Background(), &grpcPort.GetUserRequest{Id: res.Id})
	assert.NoError(t, err)
	assert.Equal(t, res.Id, usr.Id)
	assert.Equal(t, res.Name, usr.Name)
	assert.Equal(t, res.Email, usr.Email)

	cancelFunc()
}

func TestGRPCDeleteUser(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	usr, err := client.DeleteUser(context.Background(), &grpcPort.DeleteUserRequest{Id: res.Id})
	assert.NoError(t, err)
	assert.Equal(t, res.Id, usr.Id)
	assert.Equal(t, res.Name, usr.Name)
	assert.Equal(t, res.Email, usr.Email)
	_, err = client.DeleteUser(context.Background(), &grpcPort.DeleteUserRequest{Id: res.Id})
	assert.Error(t, err, ErrNotFound)

	cancelFunc()
}

func TestGRPCCreateAd(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	_, err = client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(t, err)

	cancelFunc()
}

func TestGRPCDeleteAd(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(t, err)
	ad1, err := client.DeleteAd(context.Background(), &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: ad.AuthorId})
	assert.NoError(t, err)
	assert.Equal(t, ad.Id, ad1.Id)

	cancelFunc()
}

func TestGRPCChangeAdStatus(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(t, err)
	ad1, err := client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(t, err)
	assert.Equal(t, ad.Id, ad1.Id)
	assert.Equal(t, ad1.Published, true)

	cancelFunc()
}

func TestGRPCUpdateAd(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(t, err)
	ad1, err := client.UpdateAd(context.Background(), &grpcPort.UpdateAdRequest{AdId: ad.Id, UserId: ad.AuthorId,
		Title: "world", Text: "hello"})
	assert.NoError(t, err)
	assert.Equal(t, ad.Id, ad1.Id)
	assert.Equal(t, ad1.Title, "world")
	assert.Equal(t, ad1.Text, "hello")

	cancelFunc()
}

func TestGRPCListad(t *testing.T) {

	var cancelFunc context.CancelFunc
	go func() {
		ctx, cf := context.WithCancel(context.Background())
		cancelFunc = cf
		endChan := make(chan int)
		ports.CreateServer(ctx, endChan)
		defer func() {
			<-endChan
		}()
	}()

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	res1, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(t, err)
	_, err = client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "world", Text: "hello"})
	assert.NoError(t, err)
	ad2, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res1.Id,
		Title: "hello?", Text: "world!"})
	assert.NoError(t, err)
	_, err = client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(t, err)
	ad2, err = client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad2.Id, UserId: ad2.AuthorId, Published: true})
	assert.NoError(t, err)
	ads, err := client.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_Default})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 2)
	ads, err = client.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByAuthor, Data: &grpcPort.Mode_AuthorId{AuthorId: ad2.AuthorId}})
	assert.NoError(t, err)
	assert.Len(t, ads.List, 1)

	cancelFunc()
}
