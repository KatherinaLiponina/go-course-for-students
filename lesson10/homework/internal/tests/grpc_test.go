package tests

import (
	"context"
	"log"
	"testing"

	"homework10/internal/ports"
	grpcPort "homework10/internal/ports/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcTestSuite struct {
	suite.Suite
	gsrv *grpc.Server
	cf   context.CancelFunc
	ch   chan int
	t    *testing.T
}

func (suite *GrpcTestSuite) SetupTest() {
	ctx, cf := context.WithCancel(context.Background())
	suite.cf = cf
	endChan := make(chan int)
	suite.ch = endChan
	_, suite.gsrv = ports.CreateServer(ctx, endChan)
}

func (suite *GrpcTestSuite) TearDownTest() {
	suite.cf()
	<-suite.ch
}

func TestGrpcTestSuite(t *testing.T) {
	gts := new(GrpcTestSuite)
	gts.t = t
	suite.Run(t, gts)
}

func (suite *GrpcTestSuite) TestGRPCCreateUser() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	assert.Equal(suite.t, "Oleg", res.GetName())

}

func (suite *GrpcTestSuite) TestGRPCGetUser() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	usr, err := client.GetUser(context.Background(), &grpcPort.GetUserRequest{Id: res.Id})
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, res.GetId(), usr.Id)
	assert.Equal(suite.t, res.GetName(), usr.Name)
	assert.Equal(suite.t, res.GetEmail(), usr.Email)

}

func (suite *GrpcTestSuite) TestGRPCDeleteUser() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	usr, err := client.DeleteUser(context.Background(), &grpcPort.DeleteUserRequest{Id: res.Id})
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, res.GetId(), usr.Id)
	assert.Equal(suite.t, res.GetName(), usr.Name)
	assert.Equal(suite.t, res.GetEmail(), usr.Email)
	_, err = client.DeleteUser(context.Background(), &grpcPort.DeleteUserRequest{Id: res.Id})
	assert.Error(suite.t, err, ErrNotFound)

}

func (suite *GrpcTestSuite) TestGRPCCreateAd() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	_, err = client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(suite.t, err)

}

func (suite *GrpcTestSuite) TestGRPCDeleteAd() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(suite.t, err)
	ad1, err := client.DeleteAd(context.Background(), &grpcPort.DeleteAdRequest{AdId: ad.Id, AuthorId: ad.AuthorId})
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, ad.GetId(), ad1.GetId())
}

func (suite *GrpcTestSuite) TestGRPCChangeAdStatus() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(suite.t, err)
	ad1, err := client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, ad.GetId(), ad1.GetId())
	assert.Equal(suite.t, ad1.GetPublished(), true)

}

func (suite *GrpcTestSuite) TestGRPCUpdateAd() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(suite.t, err)

	ad1, err := client.UpdateAd(context.Background(), &grpcPort.UpdateAdRequest{AdId: ad.Id, UserId: ad.AuthorId,
		Title: "world", Text: "hello"})
	assert.NoError(suite.t, err)
	assert.NotNil(suite.t, ad1.String())
	assert.Equal(suite.t, ad.GetId(), ad1.GetId())
	assert.Equal(suite.t, ad.GetAuthorId(), ad1.GetAuthorId())
	assert.Equal(suite.t, ad1.GetTitle(), "world")
	assert.Equal(suite.t, ad1.GetText(), "hello")
	assert.Greater(suite.t, ad1.GetUpdateTime().AsTime(), ad.GetUpdateTime().AsTime())
	assert.Equal(suite.t, ad1.GetCreationDate().AsTime(), ad.GetCreationDate().AsTime())

}

func (suite *GrpcTestSuite) TestGRPCListad() {

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	res1, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Name: "Oleg"})
	assert.NoError(suite.t, err, "client.GetUser")
	ad, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "hello", Text: "world"})
	assert.NoError(suite.t, err)
	_, err = client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res.Id,
		Title: "world", Text: "hello"})
	assert.NoError(suite.t, err)
	ad2, err := client.CreateAd(context.Background(), &grpcPort.CreateAdRequest{UserId: res1.Id,
		Title: "hello?", Text: "world!"})
	assert.NoError(suite.t, err)
	_, err = client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad.Id, UserId: ad.AuthorId, Published: true})
	assert.NoError(suite.t, err)
	ad2, err = client.ChangeAdStatus(context.Background(), &grpcPort.ChangeAdStatusRequest{AdId: ad2.Id, UserId: ad2.AuthorId, Published: true})
	assert.NoError(suite.t, err)
	ads, err := client.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_Default})
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.List, 2)
	ads, err = client.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByAuthor, Data: &grpcPort.Mode_AuthorId{AuthorId: ad2.AuthorId}})
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.List, 1)

}
