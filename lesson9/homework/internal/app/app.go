package app

import (
	"errors"
	"homework9/internal/ads"
	"homework9/internal/users"
	"time"

	"github.com/KatherinaLiponina/validation"
)

var ErrNotFound = errors.New("repository does not contain ad with given ID")
var ErrForbidden = errors.New("authorID does not match given ID")
var ErrBadRequest = errors.New("validation for title or text was failed")

type App interface {
	CreateAd(Title string, Text string, AuthorID int64) (*ads.Ad, error)
	ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error)
	UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error)
	GetAdByID(ID int64) (*ads.Ad, error)
	DeleteAd(ID int64, AuthorID int64) (*ads.Ad, error)

	Select() []ads.Ad
	SelectByAuthor(authorID int64) ([]ads.Ad, error)
	SelectByCreation(time time.Time) []ads.Ad
	SelectAll() []ads.Ad
	FindByTitle(Title string) []ads.Ad

	CreateUser(nickname string, email string) *users.User
	UpdateUser(ID int64, nickname string, email string) (*users.User, error)
	GetUserByID(ID int64) (*users.User, error)
	DeleteUser(ID int64) (*users.User, error)
}

type AdRepository interface {
	AppendAd(Title string, Text string, AuthorID int64) *ads.Ad
	ChangeAdStatus(ID int64, status bool)
	UpdateAd(ID int64, Text string, Title string)
	GetAdByID(ID int64) (*ads.Ad, error)
	Select(f func(ads.Ad) bool) []ads.Ad
	DeleteAd(ID int64) (*ads.Ad, error)
}

type UserRepository interface {
	AppendUser(nickname string, email string) *users.User
	UpdateUser(ID int64, nickname string, email string)
	GetUserByID(ID int64) (*users.User, error)
	DeleteUser(ID int64) (*users.User, error)
}

type app struct {
	adrepo  AdRepository
	usrrepo UserRepository
}

type validationStruct struct {
	Title string `validate:"title"`
	Text  string `validate:"text"`
}

func newValidationStruct(title string, text string) validationStruct {
	return validationStruct{Title: title, Text: text}
}

func (a *app) CreateAd(Title string, Text string, AuthorID int64) (*ads.Ad, error) {
	err := validation.Validate(newValidationStruct(Title, Text))
	if err != nil {
		return nil, ErrBadRequest
	}
	_, err = a.usrrepo.GetUserByID(AuthorID)
	if err != nil {
		return nil, ErrNotFound
	}
	return a.adrepo.AppendAd(Title, Text, AuthorID), nil
}

func (a *app) ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error) {
	ad, err := a.adrepo.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.adrepo.ChangeAdStatus(ID, status)
	return a.adrepo.GetAdByID(ID)
}

func (a *app) UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error) {
	err := validation.Validate(newValidationStruct(Title, Text))
	if err != nil {
		return nil, ErrBadRequest
	}
	_, err = a.usrrepo.GetUserByID(AuthorID)
	if err != nil {
		return nil, ErrNotFound
	}
	ad, err := a.adrepo.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.adrepo.UpdateAd(ID, Text, Title)
	return a.adrepo.GetAdByID(ID)
}

func (a *app) GetAdByID(ID int64) (*ads.Ad, error) {
	return a.adrepo.GetAdByID(ID)
}

func (a *app) Select() []ads.Ad {
	return a.adrepo.Select(func(a ads.Ad) bool { return a.Published })
}

func (a *app) SelectByAuthor(authorID int64) ([]ads.Ad, error) {
	_, err := a.usrrepo.GetUserByID(authorID)
	if err != nil {
		return nil, ErrNotFound
	}
	return a.adrepo.Select(func(a ads.Ad) bool { return a.AuthorID == authorID }), nil
}
func (a *app) SelectByCreation(time time.Time) []ads.Ad {
	return a.adrepo.Select(func(a ads.Ad) bool { return a.CreationDate.After(time) })
}

func (a *app) SelectAll() []ads.Ad {
	return a.adrepo.Select(func(a ads.Ad) bool { return true })
}

func (a *app) DeleteAd(ID int64, AuthorID int64) (*ads.Ad, error) {
	_, err := a.usrrepo.GetUserByID(AuthorID)
	if err != nil {
		return nil, ErrNotFound
	}
	ad, err := a.adrepo.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	return a.adrepo.DeleteAd(ID)
}

func (a *app) CreateUser(nickname string, email string) *users.User {
	return a.usrrepo.AppendUser(nickname, email)
}

func (a *app) UpdateUser(ID int64, nickname string, email string) (*users.User, error) {
	_, err := a.usrrepo.GetUserByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	a.usrrepo.UpdateUser(ID, nickname, email)
	return a.usrrepo.GetUserByID(ID)
}

func (a *app) FindByTitle(Title string) []ads.Ad {
	return a.adrepo.Select(func(a ads.Ad) bool {
		return a.Title == Title
	})
}

func (a *app) GetUserByID(ID int64) (*users.User, error) {
	usr, err := a.usrrepo.GetUserByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	return usr, nil
}

func (a *app) DeleteUser(ID int64) (*users.User, error) {
	_, err := a.usrrepo.GetUserByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	return a.usrrepo.DeleteUser(ID)
}

func NewApp(a AdRepository, u UserRepository) App {
	return &app{adrepo: a, usrrepo: u}
}
