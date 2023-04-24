package ads

import "time"

type Ad struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Text         string    `json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	CreationDate time.Time `json:"creation_time"`
	UpdateTime   time.Time `json:"update_time"`
}

func CreateAd(ID int64, Title string, Text string, AuthorID int64) Ad {
	current_time := time.Now().UTC()
	return Ad{ID, Title, Text, AuthorID, false, current_time, current_time}
}

func (a *Ad) ChangeAdStatus(status bool) {
	a.Published = status
}

func (a *Ad) UpdateTitle(title string) {
	a.Title = title
	a.UpdateTime = time.Now().UTC()
}

func (a *Ad) UpdateText(text string) {
	a.Text = text
	a.UpdateTime = time.Now()
}
