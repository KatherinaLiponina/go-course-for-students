package users

type User struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func CreateUser(id int64, nick string, email string) User {
	return User{id, nick, email}
}

func (u * User) UpdateNickname(n string) {
	u.Nickname = n
}

func (u * User) UpdateEmail(e string) {
	u.Email = e
}