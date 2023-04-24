package userrepo

import (
	"errors"
	"homework9/internal/app"
	"homework9/internal/users"
	"sync"
)

type repo struct {
	mtx sync.RWMutex
	index int64
	usrStorage map[int64]users.User
}

func (r * repo) AppendUser(nickname string, email string) *users.User {
	r.mtx.Lock()
	usr := users.CreateUser(r.index, nickname, email)
	r.index++
	r.usrStorage[usr.ID] = usr
	r.mtx.Unlock()
	return &usr
}

func (r * repo) UpdateUser(ID int64, nickname string, email string) {
	r.mtx.Lock()
	usr := r.usrStorage[ID]
	if len(nickname) > 0 {
		usr.UpdateNickname(nickname)
	}
	if len(email) > 0 {
		usr.UpdateEmail(email)
	}
	r.usrStorage[usr.ID] = usr
	r.mtx.Unlock()
}

func (r * repo) GetUserByID(ID int64) (*users.User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	a, ok := r.usrStorage[ID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &a, nil
}

func New() app.UserRepository {
	return &repo{index: 0, usrStorage: map[int64]users.User{}}
}