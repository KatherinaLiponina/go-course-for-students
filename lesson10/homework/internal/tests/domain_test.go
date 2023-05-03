package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()
	defer client.cancelTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()
	defer client.cancelTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()
	defer client.cancelTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

func TestCreateAdWithoutUser(t *testing.T) {
	client := getTestClient()
	defer client.cancelTestClient()
	_, err := client.createAd(1, "hello", "world")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestDeleteAdWithWrongUser(t *testing.T) {
	client := getTestClient()
	defer client.cancelTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.DeleteAd(resp.Data.ID, 1)
	assert.ErrorIs(t, err, ErrForbidden)
}
