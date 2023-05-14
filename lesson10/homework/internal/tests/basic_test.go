package tests

import (
	"context"
	"net/http"
	"testing"

	"homework10/internal/ports"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BasicTestSuite struct {
	suite.Suite
	hsrv *http.Server
	cf   context.CancelFunc
	ch   chan int
	t    *testing.T
}

func (suite *BasicTestSuite) SetupTest() {
	ctx, cf := context.WithCancel(context.Background())
	suite.cf = cf
	endChan := make(chan int)
	suite.ch = endChan
	suite.hsrv, _ = ports.CreateServer(ctx, endChan)
}

func (suite *BasicTestSuite) TearDownTest() {
	suite.cf()
	<-suite.ch
}

func TestBasicTestSuite(t *testing.T) {
	bts := new(BasicTestSuite)
	bts.t = t
	suite.Run(t, bts)
}

func (suite *BasicTestSuite) TestCreateUser() {
	client := getTestClient(suite.hsrv.Addr)
	response, err := client.createUser("Alice", "alice.doe@gmail.com")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, int64(0))
	assert.Equal(suite.t, response.Data.Nickname, "Alice")
	assert.Equal(suite.t, response.Data.Email, "alice.doe@gmail.com")
	response, err = client.createUser("Bob", "bob.doe@gmail.com")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, int64(1))
}

func (suite *BasicTestSuite) TestGetUser() {
	client := getTestClient(suite.hsrv.Addr)

	resp, err := client.createUser("David", "david.doe@gmail.com")
	assert.NoError(suite.t, err)
	response, err := client.GetUserByID(resp.Data.ID)
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, resp.Data.ID)
	assert.Equal(suite.t, response.Data.Nickname, resp.Data.Nickname)
	assert.Equal(suite.t, response.Data.Email, resp.Data.Email)
}

func (suite *BasicTestSuite) TestDeleteUser() {
	client := getTestClient(suite.hsrv.Addr)

	resp, err := client.createUser("Harry", "harry.doe@gmail.com")
	assert.NoError(suite.t, err)
	response, err := client.DeleteUser(resp.Data.ID)
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, resp.Data.ID)
	assert.Equal(suite.t, response.Data.Nickname, resp.Data.Nickname)
	assert.Equal(suite.t, response.Data.Email, resp.Data.Email)
	response, err = client.DeleteUser(resp.Data.ID)
	assert.ErrorIs(suite.t, err, ErrNotFound)
}

func (suite *BasicTestSuite) TestDeleteAd() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Carol", "carol.doe@gmail.com")
	assert.NoError(suite.t, err)
	resp2, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	resp3, err := client.DeleteAd(resp2.Data.ID, resp2.Data.AuthorID)
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, resp2.Data.ID, resp3.Data.ID)
	assert.Equal(suite.t, resp2.Data.AuthorID, resp3.Data.AuthorID)
	_, err = client.DeleteAd(resp2.Data.ID, resp2.Data.AuthorID)
	assert.ErrorIs(suite.t, err, ErrNotFound)
}

func (suite *BasicTestSuite) TestUpdateUser() {
	client := getTestClient(suite.hsrv.Addr)

	response_old, err := client.createUser("Eva", "eva.doe@gmail.com")
	assert.NoError(suite.t, err)
	response, err := client.updateUser(response_old.Data.ID, "", "new.mail@yandex.ru")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, response_old.Data.ID)
	assert.Equal(suite.t, response.Data.Nickname, "Eva")
	assert.Equal(suite.t, response.Data.Email, "new.mail@yandex.ru")

	response, err = client.updateUser(response_old.Data.ID, "NewEva", "new.mail@yandex.ru")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.ID, response_old.Data.ID)
	assert.Equal(suite.t, response.Data.Nickname, "NewEva")
}

func (suite *BasicTestSuite) TestCreateAd() {
	client := getTestClient(suite.hsrv.Addr)

	_, err := client.createUser("Franc", "franc.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr, err := client.createUser("Georg", "georg.doe@gmail.com")
	assert.NoError(suite.t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(suite.t, err)
	assert.Zero(suite.t, response.Data.ID)
	assert.Equal(suite.t, response.Data.Title, "hello")
	assert.Equal(suite.t, response.Data.Text, "world")
	assert.Equal(suite.t, response.Data.AuthorID, usr.Data.ID)
	assert.False(suite.t, response.Data.Published)
}

func (suite *BasicTestSuite) TestChangeAdStatus() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Irma", "irma.doe@gmail.com")
	assert.NoError(suite.t, err)

	response, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	response, err = client.changeAdStatus(usr.Data.ID, response.Data.ID, true)
	assert.NoError(suite.t, err)
	assert.True(suite.t, response.Data.Published)

	response, err = client.changeAdStatus(usr.Data.ID, response.Data.ID, false)
	assert.NoError(suite.t, err)
	assert.False(suite.t, response.Data.Published)

	response, err = client.changeAdStatus(usr.Data.ID, response.Data.ID, false)
	assert.NoError(suite.t, err)
	assert.False(suite.t, response.Data.Published)
}

func (suite *BasicTestSuite) TestUpdateAd() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Jane", "john.doe@gmail.com")
	assert.NoError(suite.t, err)

	response, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	response, err = client.updateAd(usr.Data.ID, response.Data.ID, "привет", "мир")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, response.Data.Title, "привет")
	assert.Equal(suite.t, response.Data.Text, "мир")
}

func (suite *BasicTestSuite) TestListAds() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Kate", "kate.doe@gmail.com")
	assert.NoError(suite.t, err)

	response, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	publishedAd, err := client.changeAdStatus(usr.Data.ID, response.Data.ID, true)
	assert.NoError(suite.t, err)

	_, err = client.createAd(usr.Data.ID, "best cat", "not for sale")
	assert.NoError(suite.t, err)

	ads, err := client.listAds(nil)
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.Data, int(1))
	assert.Equal(suite.t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(suite.t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(suite.t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(suite.t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(suite.t, ads.Data[0].Published)
}

func (suite *BasicTestSuite) TestListAdsByAuthor() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Lisa", "lisa.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Mary", "mary.doe@gmail.com")
	assert.NoError(suite.t, err)

	ad, _ := client.createAd(usr2.Data.ID, "hello", "world")
	_, err = client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	_, err = client.createAd(usr2.Data.ID, "world", "hello")
	assert.NoError(suite.t, err)

	ads, err := client.listAdsByAuthor(usr2.Data.ID)
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.Data, 2)
	assert.Equal(suite.t, ads.Data[0].AuthorID, ad.Data.AuthorID)
	assert.Equal(suite.t, ads.Data[1].AuthorID, ad.Data.AuthorID)
}

func (suite *BasicTestSuite) TestListAdsByTime() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Nansy", "nansy.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Olga", "olga.doe@gmail.com")
	assert.NoError(suite.t, err)

	ad1, err := client.createAd(usr2.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	ad2, err := client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	ads, err := client.listAdsByTime(ad1.Data.CreationTime)
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.Data, 1)
	assert.Equal(suite.t, ads.Data[0].ID, ad2.Data.ID)
	assert.Equal(suite.t, ads.Data[0].Title, ad2.Data.Title)
	assert.Equal(suite.t, ads.Data[0].Text, ad2.Data.Text)
	assert.Equal(suite.t, ads.Data[0].AuthorID, ad2.Data.AuthorID)
}

func (suite *BasicTestSuite) TestListAll() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Peter", "peter.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Rosa", "rosa.doe@gmail.com")
	assert.NoError(suite.t, err)

	_, err = client.createAd(usr2.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	_, err = client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	ads, err := client.listAll()
	assert.NoError(suite.t, err)
	assert.Len(suite.t, ads.Data, 2)
}

func (suite *BasicTestSuite) TestGetAdById() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Sally", "sally.doe@gmail.com")
	assert.NoError(suite.t, err)

	ad, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	adAgain, err := client.getAd(ad.Data.ID)
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, ad.Data.ID, adAgain.Data.ID)
}

func (suite *BasicTestSuite) TestFindByName() {
	client := getTestClient(suite.hsrv.Addr)

	usr, err := client.createUser("Tuomas", "tuomas.doe@gmail.com")
	assert.NoError(suite.t, err)

	ad, err := client.createAd(usr.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	_, err = client.createAd(usr.Data.ID, "hello", "mir")
	assert.NoError(suite.t, err)
	_, err = client.createAd(usr.Data.ID, "hallo", "welt")
	assert.NoError(suite.t, err)
	arr, err := client.findAdByTitle(ad.Data.Title)
	assert.NoError(suite.t, err)
	assert.Len(suite.t, arr.Data, 2)
	assert.Equal(suite.t, arr.Data[0].Title, ad.Data.Title)
	assert.Equal(suite.t, arr.Data[1].Title, ad.Data.Title)
}
