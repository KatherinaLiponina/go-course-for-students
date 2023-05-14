package tests

import (
	"context"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/ports"
	"homework10/internal/tests/mocks"
	"homework10/internal/users"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAppWithRepoMock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	adrepo := mocks.NewMockAdRepository(mockCtrl)
	usrrepo := mocks.NewMockUserRepository(mockCtrl)

	testad := &ads.Ad{Title: "Title", Text: "Text", AuthorID: 0}
	adrepo.EXPECT().AppendAd(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testad)
	adrepo.EXPECT().ChangeAdStatus(gomock.Any(), gomock.Any()).AnyTimes()
	adrepo.EXPECT().GetAdByID(gomock.Any()).AnyTimes().Return(testad, nil)
	adrepo.EXPECT().DeleteAd(gomock.Any()).Return(testad, nil).AnyTimes()
	adrepo.EXPECT().UpdateAd(gomock.Any(), gomock.Any(), gomock.Any())
	adrepo.EXPECT().Select(gomock.Any()).AnyTimes().Return([]ads.Ad{*testad})

	testusr := &users.User{ID: 0, Nickname: "Test Subject", Email: "glados@aparture.com"}
	usrrepo.EXPECT().AppendUser(gomock.Any(), gomock.Any()).AnyTimes().Return(testusr)
	usrrepo.EXPECT().DeleteUser(gomock.Any()).Return(testusr, nil)
	usrrepo.EXPECT().GetUserByID(gomock.Any()).Return(testusr, nil).AnyTimes()
	usrrepo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any())

	a := app.NewApp(adrepo, usrrepo)
	usr := a.CreateUser("Chell", "chell@mail.org")
	assert.Equal(t, usr.ID, testusr.ID)
	assert.Equal(t, usr.Nickname, testusr.Nickname)
	assert.Equal(t, usr.Email, testusr.Email)

	usr, err := a.DeleteUser(usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, usr.ID, testusr.ID)

	usr, err = a.GetUserByID(usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, usr.ID, testusr.ID)

	usr, err = a.UpdateUser(usr.ID, "Jane", "Email")
	assert.NoError(t, err)
	assert.Equal(t, usr.ID, testusr.ID)
	assert.Equal(t, usr.Nickname, testusr.Nickname)
	assert.Equal(t, usr.Email, testusr.Email)

	ad, err := a.CreateAd("NotTitle", "NotText", usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.ID, testad.ID)
	assert.Equal(t, ad.Title, testad.Title)
	assert.Equal(t, ad.Text, testad.Text)

	ad, err = a.GetAdByID(ad.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.ID, testad.ID)

	adarr := a.Select()
	assert.Len(t, adarr, 1)

	_, err = a.DeleteAd(testad.ID, 9)
	assert.Error(t, err)

	ad, err = a.DeleteAd(testad.ID, usr.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.ID, testad.ID)

	ad, err = a.UpdateAd(ad.ID, ad.AuthorID, "a", "b")
	assert.NoError(t, err)
	assert.Equal(t, ad.ID, testad.ID)

	ad, err = a.ChangeAdStatus(ad.ID, ad.AuthorID, true)
	assert.NoError(t, err)
	assert.Equal(t, ad.ID, testad.ID)
}

func TestHandlerWithAppMock(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	testad := &ads.Ad{Title: "Title", Text: "Text", AuthorID: 0}
	testusr := &users.User{ID: 0, Nickname: "Test Subject", Email: "glados@aparture.com"}

	appmock := mocks.NewMockApp(mockCtrl)
	appmock.EXPECT().ChangeAdStatus(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(ID int64, AuthorID int64, status bool) (*ads.Ad, error) {
			if AuthorID != 0 {
				return &ads.Ad{}, ErrForbidden
			}
			testad.Published = status
			return testad, nil
		})
	appmock.EXPECT().CreateAd(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testad, nil)
	appmock.EXPECT().CreateUser(gomock.Any(), gomock.Any()).AnyTimes().Return(testusr)
	appmock.EXPECT().DeleteAd(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(ID int64, AuthorID int64) (*ads.Ad, error) {
		if AuthorID != 0 {
			return &ads.Ad{}, ErrForbidden
		}
		return testad, nil
	})
	appmock.EXPECT().DeleteUser(gomock.Any()).AnyTimes().DoAndReturn(func(ID int64) (*users.User, error) {
		if ID != testusr.ID {
			return &users.User{}, ErrNotFound
		}
		return testusr, nil
	})
	appmock.EXPECT().GetAdByID(gomock.Any()).AnyTimes().Return(testad, nil)
	appmock.EXPECT().GetUserByID(gomock.Any()).AnyTimes().Return(testusr, nil)
	appmock.EXPECT().UpdateAd(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testad, nil)
	appmock.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(testusr, nil)

	appmock.EXPECT().Select().AnyTimes().Return([]ads.Ad{*testad})
	appmock.EXPECT().SelectAll().AnyTimes().Return([]ads.Ad{*testad})
	appmock.EXPECT().SelectByAuthor(gomock.Any()).AnyTimes().Return([]ads.Ad{*testad}, nil)
	appmock.EXPECT().SelectByCreation(gomock.Any()).AnyTimes().Return([]ads.Ad{*testad})

	ctx, cf := context.WithCancel(context.Background())
	endChan := make(chan int)
	hsrv, _ := ports.CreateServerWithExternalApp(ctx, endChan, appmock)

	client := getTestClient(hsrv.Addr)
	resp, err := client.createUser("Alice", "alice.doe@gmail.com")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, testusr.ID)
	assert.Equal(t, resp.Data.Nickname, testusr.Nickname)

	resp, err = client.GetUserByID(12)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, testusr.ID)

	resp, err = client.DeleteUser(0)
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, testusr.ID)

	ad, err := client.createAd(0, "some", "some")
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, testad.ID)
	assert.Equal(t, ad.Data.Title, testad.Title)

	ad, err = client.changeAdStatus(ad.Data.AuthorID, ad.Data.ID, true)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.Published, true)

	ad, err = client.updateAd(0, 0, "title", "text")
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, testad.ID)
	assert.Equal(t, ad.Data.Title, testad.Title)

	ad, err = client.getAd(0)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, testad.ID)

	ad, err = client.DeleteAd(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, testad.ID)

	resp, err = client.updateUser(0, "Free", "Now")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, testusr.ID)
	assert.Equal(t, resp.Data.Nickname, testusr.Nickname)

	adsArr, err := client.listAds(nil)
	assert.NoError(t, err)
	assert.Len(t, adsArr.Data, 1)

	adsArr, err = client.listAll()
	assert.NoError(t, err)
	assert.Len(t, adsArr.Data, 1)

	adsArr, err = client.listAdsByAuthor(0)
	assert.NoError(t, err)
	assert.Len(t, adsArr.Data, 1)

	adsArr, err = client.listAdsByTime(time.Now())
	assert.NoError(t, err)
	assert.Len(t, adsArr.Data, 1)

	cf()
	<-endChan
}
