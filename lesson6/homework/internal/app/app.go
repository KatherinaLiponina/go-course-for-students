package app

import (
	"errors"
	"homework6/internal/ads"
)

var ErrNotFound = errors.New("Repository does not contain ad with given ID")
var ErrForbidden = errors.New("AuthorID does not match given ID")

type App interface {
	CreateAd(Title string, Text string, AuthorID int64) (*ads.Ad, error)
	ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error)
	UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error)
}

type Repository interface {
	AppendAd(ad ads.Ad)
	ChangeAdStatus(ID int64, status bool)
	UpdateAd(ID int64, Text string, Title string)
	GetAdByID(ID int64) (*ads.Ad, error)
}

type app struct {
	index int64
	repository Repository
}

func (a * app) CreateAd(Title string, Text string, AuthorID int64) (*ads.Ad, error) {
	ad := ads.CreateAd(a.index, Title, Text, AuthorID)
	a.index++
	a.repository.AppendAd(ad)
	return &ad, nil
}

func (a * app) ChangeAdStatus(ID int64, AuthorID int64, status bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.repository.ChangeAdStatus(ID, status)
	return a.repository.GetAdByID(ID);
}

func (a * app) UpdateAd(ID int64, AuthorID int64, Title string, Text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ID)
	if err != nil {
		return nil, ErrNotFound
	}
	if ad.AuthorID != AuthorID {
		return nil, ErrForbidden
	}
	a.repository.UpdateAd(ID, Text, Title)
	return a.repository.GetAdByID(ID);
}

func NewApp(repo Repository) App {
	return &app{repository: repo, index: 0}
}
