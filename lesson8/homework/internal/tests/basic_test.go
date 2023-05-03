package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t * testing.T) {
	client := getTestClient()

	response, err := client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(0))
	assert.Equal(t, response.Data.Nickname, "Jane")
	assert.Equal(t, response.Data.Email, "jane.doe@gmail.com")
	response, err = client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(1))
}

func TestUpdateUser(t * testing.T) {
	client := getTestClient()
	response, err := client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)
	response, err = client.updateUser(response.Data.ID, "", "new.mail@yandex.ru")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, int64(0))
	assert.Equal(t, response.Data.Nickname, "Jane")
	assert.Equal(t, response.Data.Email, "new.mail@yandex.ru")
}

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(1))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(1, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	response, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(1, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds(nil)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestListAdsByAuthor(t * testing.T) {
	client := getTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	ad, _ := client.createAd(1, "hello", "world")
	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(1, "world", "hello")
	assert.NoError(t, err)

	ads, err := client.listAdsByAuthor(1)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	assert.Equal(t, ads.Data[0].AuthorID, ad.Data.AuthorID)
	assert.Equal(t, ads.Data[1].AuthorID, ad.Data.AuthorID)
}

func TestListAdsByTime(t * testing.T) {
	client := getTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	ad1, err := client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	ad2, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ads, err := client.listAdsByTime(ad1.Data.CreationTime)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, ad2.Data.ID)
	assert.Equal(t, ads.Data[0].Title, ad2.Data.Title)
	assert.Equal(t, ads.Data[0].Text, ad2.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, ad2.Data.AuthorID)
}

func TestListAll(t * testing.T) {
	client := getTestClient()
	_, err := client.createUser("John", "john.doe@gmail.com")
	assert.NoError(t, err)
	_, err = client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	_, err = client.createAd(1, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	ads, err := client.listAll()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

func TestGetAdById(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	ad, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	adAgain, err := client.getAd(ad.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, adAgain.Data.ID)
}

func TestFindByName(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser("Jane", "jane.doe@gmail.com")
	assert.NoError(t, err)

	ad, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hello", "mir")
	assert.NoError(t, err)
	_, err = client.createAd(0, "hallo", "welt")
	assert.NoError(t, err)
	arr, err := client.findAdByTitle(ad.Data.Title)
	assert.NoError(t, err)
	assert.Len(t, arr.Data, 2)
	assert.Equal(t, arr.Data[0].Title, ad.Data.Title)
	assert.Equal(t, arr.Data[1].Title, ad.Data.Title)
}
