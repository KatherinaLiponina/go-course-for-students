package adrepo

import (
	"errors"
	"homework10/internal/ads"
	"homework10/internal/app"
	"sync"
)

type repo struct {
	mtx       sync.RWMutex
	index     int64
	adStorage map[int64]ads.Ad
}

func (r *repo) AppendAd(Title string, Text string, AuthorID int64) *ads.Ad {
	r.mtx.Lock()
	ad := ads.CreateAd(r.index, Title, Text, AuthorID)
	r.index++
	r.adStorage[ad.ID] = ad
	r.mtx.Unlock()
	return &ad
}

func (r *repo) ChangeAdStatus(ID int64, status bool) {
	r.mtx.Lock()
	ad := r.adStorage[ID]
	ad.ChangeAdStatus(status)
	r.adStorage[ad.ID] = ad
	r.mtx.Unlock()
}

func (r *repo) UpdateAd(ID int64, Text string, Title string) {
	r.mtx.Lock()
	ad := r.adStorage[ID]
	if len(Text) > 0 {
		ad.UpdateText(Text)
	}
	if len(Title) > 0 {
		ad.UpdateTitle(Title)
	}
	r.adStorage[ad.ID] = ad
	r.mtx.Unlock()
}

func (r *repo) GetAdByID(ID int64) (*ads.Ad, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	a, ok := r.adStorage[ID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &a, nil
}

func (r *repo) Select(f func(ads.Ad) bool) []ads.Ad {
	r.mtx.RLock()
	resultArray := make([]ads.Ad, 0)
	for _, v := range r.adStorage {
		if f(v) {
			resultArray = append(resultArray, v)
		}
	}
	r.mtx.RUnlock()
	return resultArray
}

func (r *repo) DeleteAd(ID int64) (*ads.Ad, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	a, ok := r.adStorage[ID]
	if !ok {
		return nil, errors.New("not found")
	}
	delete(r.adStorage, a.ID)
	return &a, nil
}

func New() app.AdRepository {
	return &repo{index: 0, adStorage: map[int64]ads.Ad{}}
}
