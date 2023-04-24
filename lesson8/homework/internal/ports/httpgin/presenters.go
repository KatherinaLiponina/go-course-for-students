package httpgin

import (
	"homework8/internal/ads"
	"homework8/internal/users"
	"time"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type selectAdRequest struct {
	ByAuthor     bool      `json:"by_author"`
	AuthorID     int64     `json:"author_id"`
	ByCreation   bool      `json:"by_creation"`
	CreationTime time.Time `json:"creation_time"`
	All          bool      `json:"all"`
}

type createOrUpdateUser struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	Data users.User `json:"data"`
}

type adResponse struct {
	Data ads.Ad `json:"data"`
}

type adsResponse struct {
	Data []ads.Ad `json:"data"`
}