package adrepo

import (
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repo struct {
	index int64
	adStorage map[int64]ads.Ad
}

func (r * repo) AppendAd(Title string, Text string, AuthorID int64) *ads.Ad {
	ad := ads.CreateAd(r.index, Title, Text, AuthorID)
	r.index++
	r.adStorage[ad.ID] = ad
	return &ad
}

func (r * repo) ChangeAdStatus(ID int64, status bool) {
	ad, _ := r.GetAdByID(ID)
	ad.ChangeAdStatus(status)
	r.adStorage[ad.ID] = *ad
}

func (r * repo) UpdateAd(ID int64, Text string, Title string) {
	ad, _ := r.GetAdByID(ID)
	if len(Text) > 0 {
		ad.UpdateText(Text)
	}
	if len(Title) > 0 {
		ad.UpdateTitle(Title)
	}
	r.adStorage[ad.ID] = *ad
}

func (r * repo) GetAdByID(ID int64) (*ads.Ad, error) {
	a, ok := r.adStorage[ID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &a, nil
}

func New() app.Repository {
	return &repo{0, map[int64]ads.Ad{}}
}
