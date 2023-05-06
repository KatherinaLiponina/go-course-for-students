package tests

import (
	"context"
	"homework10/internal/ports"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DomainTestSuite struct {
	suite.Suite
	hsrv *http.Server
	cf   context.CancelFunc
	ch   chan int
	t    *testing.T
}

func (suite *DomainTestSuite) SetupTest() {
	ctx, cf := context.WithCancel(context.Background())
	suite.cf = cf
	endChan := make(chan int)
	suite.ch = endChan
	suite.hsrv, _ = ports.CreateServer(ctx, endChan)
}

func (suite *DomainTestSuite) TearDownTest() {
	suite.cf()
	<-suite.ch
}

func TestDomainTestSuite(t *testing.T) {
	dts := new(DomainTestSuite)
	dts.t = t
	suite.Run(t, dts)
}

func (suite *DomainTestSuite) TestChangeStatusAdOfAnotherUser() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Uma", "uma.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Victor", "jane.doe@gmail.com")
	assert.NoError(suite.t, err)

	resp, err := client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.changeAdStatus(usr2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(suite.t, err, ErrForbidden)
}

func (suite *DomainTestSuite) TestUpdateAdOfAnotherUser() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Wolfgang", "wolfgang.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Xandria", "xandria.doe@gmail.com")
	assert.NoError(suite.t, err)

	resp, err := client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.updateAd(usr2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(suite.t, err, ErrForbidden)
}

func (suite *DomainTestSuite) TestCreateAd_ID() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Yan", "yan.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Zigfrid", "zigfrid.doe@gmail.com")
	assert.NoError(suite.t, err)

	resp, err := client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, resp.Data.ID, int64(0))

	resp, err = client.createAd(usr2.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, resp.Data.ID, int64(1))

	resp, err = client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)
	assert.Equal(suite.t, resp.Data.ID, int64(2))
}

func (suite *DomainTestSuite) TestCreateAdWithoutUser() {
	client := getTestClient(suite.hsrv.Addr)

	_, err := client.createAd(124, "hello", "world")
	assert.ErrorIs(suite.t, err, ErrNotFound)
}

func (suite *DomainTestSuite) TestDeleteAdWithWrongUser() {
	client := getTestClient(suite.hsrv.Addr)

	usr1, err := client.createUser("Quark", "quark.doe@gmail.com")
	assert.NoError(suite.t, err)
	usr2, err := client.createUser("Quel", "quel.doe@gmail.com")
	assert.NoError(suite.t, err)

	resp, err := client.createAd(usr1.Data.ID, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.DeleteAd(resp.Data.ID, usr2.Data.ID)
	assert.ErrorIs(suite.t, err, ErrForbidden)
}
