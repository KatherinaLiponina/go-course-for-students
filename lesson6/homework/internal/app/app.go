package app

import (
	"github.com/KatherinaLiponina/validation"

	"errors"
	"homework6/internal/ads"
)

var ErrNotFound = errors.New("repository does not contain ad with given ID")
var ErrForbidden = errors.New("authorID does not match given ID")
var ErrBadRequest = errors.New("validation for title or text was failed")

type App interface {
	CreateAd(Title string, Text string, AuthorID int64) (*ads.Ad, error)
	ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error)
	UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error)
}

type Repository interface {
	AppendAd(Title string, Text string, AuthorID int64) *ads.Ad
	ChangeAdStatus(ID int64, status bool)
	UpdateAd(ID int64, Text string, Title string)
	GetAdByID(ID int64) (*ads.Ad, error)
}

type app struct {
	repository Repository
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
	return a.repository.AppendAd(Title, Text, AuthorID), nil
}

func (a *app) ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.repository.ChangeAdStatus(ID, status)
	return a.repository.GetAdByID(ID)
}

func (a *app) UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error) {
	err := validation.Validate(newValidationStruct(Title, Text))
	if err != nil {
		return nil, ErrBadRequest
	}
	ad, err := a.repository.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.repository.UpdateAd(ID, Text, Title)
	return a.repository.GetAdByID(ID)
}

func NewApp(repo Repository) App {
	return &app{repository: repo}
}
