package lyrics

type SearchResult struct {
	Meta     Meta     `json:"meta"`
	Response Response `json:"response"`
}
type Meta struct {
	Status int `json:"status"`
}
type Stats struct {
	UnreviewedAnnotations int  `json:"unreviewed_annotations"`
	Hot                   bool `json:"hot"`
	Pageviews             int  `json:"pageviews"`
}
type PrimaryArtist struct {
	APIPath        string `json:"api_path"`
	HeaderImageURL string `json:"header_image_url"`
	ID             int    `json:"id"`
	ImageURL       string `json:"image_url"`
	IsMemeVerified bool   `json:"is_meme_verified"`
	IsVerified     bool   `json:"is_verified"`
	Name           string `json:"name"`
	URL            string `json:"url"`
	Iq             int    `json:"iq"`
}
type Result struct {
	AnnotationCount          int            `json:"annotation_count"`
	APIPath                  string         `json:"api_path"`
	FullTitle                string         `json:"full_title"`
	HeaderImageThumbnailURL  string         `json:"header_image_thumbnail_url"`
	HeaderImageURL           string         `json:"header_image_url"`
	ID                       int            `json:"id"`
	LyricsOwnerID            int            `json:"lyrics_owner_id"`
	LyricsState              string         `json:"lyrics_state"`
	Path                     string         `json:"path"`
	PyongsCount              int            `json:"pyongs_count"`
	SongArtImageThumbnailURL string         `json:"song_art_image_thumbnail_url"`
	SongArtImageURL          string         `json:"song_art_image_url"`
	Stats                    *Stats         `json:"stats"`
	Title                    string         `json:"title"`
	TitleWithFeatured        string         `json:"title_with_featured"`
	URL                      string         `json:"url"`
	SongArtPrimaryColor      string         `json:"song_art_primary_color"`
	SongArtSecondaryColor    string         `json:"song_art_secondary_color"`
	SongArtTextColor         string         `json:"song_art_text_color"`
	PrimaryArtist            *PrimaryArtist `json:"primary_artist"`
}
type Hits struct {
	Highlights []interface{} `json:"highlights"`
	Index      string        `json:"index"`
	Type       string        `json:"type"`
	Result     *Result       `json:"result"`
}
type Response struct {
	Hits []Hits `json:"hits"`
}
