package utilities

import (
	"fmt"
	"gitlab.com/toby3d/telegraph"
)

func PostLyrics(track string, artist string, lyrics string, tg *telegraph.Account) (string, error) {
	html := fmt.Sprintf(`
<body>
%s
</body>
`, lyrics)
	data, err := telegraph.ContentFormat(html)
	if err != nil {
		return "", err
	}
	page, err := tg.CreatePage(telegraph.Page{Title: fmt.Sprintf("%s - %s [Lyrics]", artist, track),
		Content: data, AuthorName: tg.AuthorName, AuthorURL: tg.AuthorURL}, true)
	if err != nil {
		return "", err
	}
	return page.URL, nil
}
