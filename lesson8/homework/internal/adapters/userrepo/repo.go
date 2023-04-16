package userrepo

import (
	"errors"
	"homework8/internal/app"
	"homework8/internal/users"
)

type repo struct {
	index int64
	adStorage map[int64]users.User
}

func (r * repo) AppendUser(nickname string, email string) *users.User {
	usr := users.CreateUser(r.index, nickname, email)
	r.index++
	r.adStorage[usr.ID] = usr
	return &usr
}

func (r * repo) UpdateUser(ID int64, nickname string, email string) {
	usr, _ := r.GetUserByID(ID)
	if len(nickname) > 0 {
		usr.UpdateNickname(nickname)
	}
	if len(email) > 0 {
		usr.UpdateEmail(email)
	}
	r.adStorage[usr.ID] = *usr
}

func (r * repo) GetUserByID(ID int64) (*users.User, error) {
	a, ok := r.adStorage[ID]
	if !ok {
		return nil, errors.New("not found")
	}
	return &a, nil
}

func New() app.UserRepository {
	return &repo{0, map[int64]users.User{}}
}