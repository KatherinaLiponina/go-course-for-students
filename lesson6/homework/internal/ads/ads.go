package ads

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	Published bool
}

func CreateAd(ID int64, Title string, Text string, AuthorID int64) Ad {
	return Ad{ID, Title, Text, AuthorID, false}
}

func (a * Ad) ChangeAdStatus(status bool) {
	a.Published = status
}

func (a * Ad) UpdateTitle(title string) {
	a.Title = title
}

func (a * Ad) UpdateText(text string) {
	a.Text = text
}