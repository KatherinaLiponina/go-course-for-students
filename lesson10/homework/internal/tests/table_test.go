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

	var ad_data = []struct {
		id int64
		title string
		text string
	}{
		{anakin.Data.ID, "Letter to Luke", "I am your father!"},
		{anakin.Data.ID, "Letter to Lea", "I find your lack of faith disturbing"},
		{lea.Data.ID, "No", "I'd just as soon kiss a Wookiee"},
		{luke.Data.ID, "Hello from Yoda", "Do, or do not. There is no try."},
		{lea.Data.ID, "May the Force be with you", "Star wars day"},
		{luke.Data.ID, "May the Force be with you", "Star wars day"},
	}

	var ads []adResponse
	for _, d := range(ad_data) {
		ad, err := httpclient.createAd(d.id, d.title, d.text)
		assert.NoError(t, err)
		_, err = httpclient.changeAdStatus(ad.Data.AuthorID, ad.Data.ID, true)
		assert.NoError(t, err)
		ads = append(ads, ad)
	}

	_, err = httpclient.changeAdStatus(ads[3].Data.AuthorID, ads[3].Data.ID, false)
	assert.NoError(t, err)
	_, err = httpclient.changeAdStatus(ads[4].Data.AuthorID, ads[4].Data.ID, false)
	assert.NoError(t, err)

	t.Run("Default list using http", func(t *testing.T) {
		ads_resp, err := httpclient.listAds(nil)
		assert.NoError(t, err)
		assert.Equal(t, len(ads_resp.Data), 4)
	})
	
	t.Run("List all using http", func(t *testing.T) {
		ads_resp, err := httpclient.listAll()
		assert.NoError(t, err)
		assert.Equal(t, len(ads_resp.Data), 6)
	})
	
	t.Run("Default list using grpc", func(t *testing.T) {
		ads_resp_grpc, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_Default})
		assert.NoError(t, err)
		assert.Len(t, ads_resp_grpc.List, 4)
	})
	
	t.Run("List all using grpc", func(t *testing.T) {
		ads_resp_grpc, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_All})
		assert.NoError(t, err)
		assert.Len(t, ads_resp_grpc.List, 6)
	})
	

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
		{ads[0].Data.CreationTime, 5},
		{ads[2].Data.CreationTime, 3},
		{ads[4].Data.CreationTime, 1},
	}
	var byTitleSearchTest = []struct {
		title  string
		expLen int
	}{
		{"Letter", 2},
		{"Yoda", 1},
		{"Death", 0},
	}

	t.Run("Search by author using http", func(t *testing.T) {
		for _, tt := range byAuthorSearchTest {
			res, err := httpclient.listAdsByAuthor(tt.author)
			assert.NoError(t, err)
			assert.Len(t, res.Data, tt.expLen)
		}
	})
	
	t.Run("Search by date using http", func(t *testing.T) {
		for _, tt := range byDateSearchTest {
			res, err := httpclient.listAdsByTime(tt.creationTime)
			assert.NoError(t, err)
			assert.Len(t, res.Data, tt.expLen)
		}
	})
	
	t.Run("Search by title using http", func(t *testing.T) {
		for _, tt := range byTitleSearchTest {
			res, err := httpclient.findAdByTitle(tt.title)
			assert.NoError(t, err)
			assert.Len(t, res.Data, tt.expLen)
		}
	})
	
	t.Run("Search by author using grpc", func(t *testing.T) {
		for _, tt := range byAuthorSearchTest {
			res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByAuthor, Data: &grpcPort.Mode_AuthorId{AuthorId: tt.author}})
			assert.NoError(t, err)
			assert.Len(t, res.List, tt.expLen)
		}
	})
	
	t.Run("Search by date using grpc", func(t *testing.T) {
		for _, tt := range byDateSearchTest {
			res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByCreation, Data: &grpcPort.Mode_Time{Time: timestamppb.New(tt.creationTime)}})
			assert.NoError(t, err)
			assert.Len(t, res.List, tt.expLen)
		}
	})
	
	t.Run("Search by title using grpc", func(t *testing.T) {
		for _, tt := range byTitleSearchTest {
			res, err := grpcclient.ListAds(context.Background(), &grpcPort.Mode{Mode: grpcPort.ModeType_ByTitle, Data: &grpcPort.Mode_Title{Title: tt.title}})
			assert.NoError(t, err)
			assert.Len(t, res.List, tt.expLen)
		}
	})

	cf()
	<-endChan
}
