package adrepo

import (
	"errors"
	"homework6/internal/ads"
	"homework6/internal/app"
)

type repo struct {
	adStorage map[int64]ads.Ad
}

func (r * repo) AppendAd(ad ads.Ad) {
	r.adStorage[ad.ID] = ad
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

type repository struct {
	*repo
}

func New() app.Repository {
	return repository{&repo{map[int64]ads.Ad{}}}
}
