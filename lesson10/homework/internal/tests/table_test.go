package tests

import (
	"context"
	"homework10/internal/ports"
	grpcPort "homework10/internal/ports/grpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestServerUsingTable(t *testing.T) {
	ctx, cf := context.WithCancel(context.Background())
	endChan := make(chan int)
	hsrv, _ := ports.CreateServer(ctx, endChan)
	httpclient := getTestClient(hsrv.Addr)

	conn, err := grpc.DialContext(context.Background(),
		"localhost:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	grpcclient := grpcPort.NewAdServiceClient(conn)

	anakin, err := httpclient.createUser("Anakin", "anakin.skywalker@mail.com")
	assert.NoError(t, err)
	luke, err := httpclient.createUser("Luke", "luke.skywalker@mail.com")
	assert.NoError(t, err)
	lea, err := httpclient.createUser("Lea", "lea.skywalker@mail.com")
	assert.NoError(t, err)

	ad1, err := httpclient.createAd(anakin.Data.ID, "Letter to Luke", "I am your father!")
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ad1.Data.AuthorID, ad1.Data.ID, true)
	assert.NoError(t, err)
	ad2, err := httpclient.createAd(anakin.Data.ID, "Letter to Lea", "I find your lack of faith disturbing")
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ad2.Data.AuthorID, ad2.Data.ID, true)
	assert.NoError(t, err)
	ad3, err := httpclient.createAd(lea.Data.ID, "No", "I'd just as soon kiss a Wookiee")
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ad3.Data.AuthorID, ad3.Data.ID, true)
	assert.NoError(t, err)
	ad4, err := httpclient.createAd(luke.Data.ID, "Hello from Yoda", "Do, or do not. There is no try.")
	assert.NoError(t, err)
	ad5, err := httpclient.createAd(lea.Data.ID, "May the Force be with you", "Star wars day")
	assert.NoError(t, err)
	ad6, err := httpclient.createAd(luke.Data.ID, "May the Force be with you", "Star wars day")
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ad6.Data.AuthorID, ad6.Data.ID, true)
	assert.NoError(t, err)

	ads, err := httpclient.listAds(nil)
	assert.NoError(t, err)
	assert.Equal(t, len(ads.Data), 4)

	ads, err = httpclient.listAll()
	assert.NoError(t, err)
	assert.Equal(t, len(ads.Data), 6)

	adsg, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_Default})
	assert.NoError(t, err)
	assert.Len(t, adsg.List, 4)

	adsg, err = grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_All})
	assert.NoError(t, err)
	assert.Len(t, adsg.List, 6)

	_, err = httpclient.changeAdStatus(ad4.Data.AuthorID, ad4.Data.ID, true)
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ad5.Data.AuthorID, ad5.Data.ID, true)
	assert.NoError(t, err)

	var byAuthorSearchTest = []struct {
		author int64
		expLen int
	}{
		{anakin.Data.ID, 2},
		{lea.Data.ID, 2},
		{luke.Data.ID, 2},
	}
	var byDateSearchTest = []struct {
		creationTime time.Time
		expLen       int
	}{
		{ad1.Data.CreationTime, 5},
		{ad3.Data.CreationTime, 3},
		{ad5.Data.CreationTime, 1},
	}
	var byTitleSearchTest = []struct {
		title  string
		expLen int
	}{
		{"Letter", 2},
		{"Yoda", 1},
		{"Death", 0},
	}

	for _, tt := range byAuthorSearchTest {
		res, err := httpclient.listAdsByAuthor(tt.author)
		assert.NoError(t, err)
		assert.Len(t, res.Data, tt.expLen)
	}
	for _, tt := range byDateSearchTest {
		res, err := httpclient.listAdsByTime(tt.creationTime)
		assert.NoError(t, err)
		assert.Len(t, res.Data, tt.expLen)
	}
	for _, tt := range byTitleSearchTest {
		res, err := httpclient.findAdByTitle(tt.title)
		assert.NoError(t, err)
		assert.Len(t, res.Data, tt.expLen)
	}

	for _, tt := range byAuthorSearchTest {
		res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByAuthor, Data: &grpcPort.Mode_AuthorId{AuthorId: tt.author}})
		assert.NoError(t, err)
		assert.Len(t, res.List, tt.expLen)
	}
	for _, tt := range byDateSearchTest {
		res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByCreation, Data: &grpcPort.Mode_Time{Time: timestamppb.New(tt.creationTime)}})
		assert.NoError(t, err)
		assert.Len(t, res.List, tt.expLen)
	}
	for _, tt := range byTitleSearchTest {
		res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByTitle, Data: &grpcPort.Mode_Title{Title: tt.title}})
		assert.NoError(t, err)
		assert.Len(t, res.List, tt.expLen)
	}

	cf()
	<-endChan
}
