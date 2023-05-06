package tests

import (
	"context"
	"homework10/internal/ports"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ValidatonTestSuite struct {
	suite.Suite
	hsrv *http.Server
	cf   context.CancelFunc
	ch   chan int
	t    *testing.T
}

func (suite *ValidatonTestSuite) SetupTest() {
	ctx, cf := context.WithCancel(context.Background())
	suite.cf = cf
	endChan := make(chan int)
	suite.ch = endChan
	suite.hsrv, _ = ports.CreateServer(ctx, endChan)
	client := getTestClient(suite.hsrv.Addr)
	client.createUser("Admin", "admin@powerful.com")
}

func (suite *ValidatonTestSuite) TearDownTest() {
	suite.cf()
	<-suite.ch
}

func TestValidatonTestSuite(t *testing.T) {
	vts := new(ValidatonTestSuite)
	vts.t = t
	suite.Run(t, vts)
}

func (suite *ValidatonTestSuite) TestCreateAd_EmptyTitle() {
	client := getTestClient(suite.hsrv.Addr)

	_, err := client.createAd(0, "", "world")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestCreateAd_TooLongTitle() {
	client := getTestClient(suite.hsrv.Addr)

	title := strings.Repeat("a", 101)

	_, err := client.createAd(0, title, "world")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestCreateAd_EmptyText() {
	client := getTestClient(suite.hsrv.Addr)

	_, err := client.createAd(0, "title", "")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestCreateAd_TooLongText() {
	client := getTestClient(suite.hsrv.Addr)

	text := strings.Repeat("a", 501)

	_, err := client.createAd(123, "title", text)
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestUpdateAd_EmptyTitle() {
	client := getTestClient(suite.hsrv.Addr)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.updateAd(0, resp.Data.ID, "", "new_world")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestUpdateAd_TooLongTitle() {
	client := getTestClient(suite.hsrv.Addr)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(suite.t, err)

	title := strings.Repeat("a", 101)

	_, err = client.updateAd(0, resp.Data.ID, title, "world")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestUpdateAd_EmptyText() {
	client := getTestClient(suite.hsrv.Addr)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.updateAd(0, resp.Data.ID, "title", "")
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}

func (suite *ValidatonTestSuite) TestUpdateAd_TooLongText() {
	client := getTestClient(suite.hsrv.Addr)

	text := strings.Repeat("a", 501)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(suite.t, err)

	_, err = client.updateAd(0, resp.Data.ID, "title", text)
	assert.ErrorIs(suite.t, err, ErrBadRequest)
}
