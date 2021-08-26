package last_fm

type GetRecentTracks struct {
	Recenttracks *Recenttracks `json:"recenttracks,omitempty"`
	Error        int           `json:"error,omitempty"`
	Message      string        `json:"message,omitempty"`
}

type GetTopArtistsResponse struct {
	Artists Artists `json:"artists"`
}

type Artist struct {
	Name       string  `json:"name,omitempty"`
	Playcount  string  `json:"playcount,omitempty"`
	Listeners  string  `json:"listeners,omitempty"`
	Mbid       string  `json:"mbid,omitempty"`
	URL        string  `json:"url,omitempty"`
	Streamable string  `json:"streamable,omitempty"`
	Image      []Image `json:"image,omitempty"`
}

type Attr struct {
	Page       string `json:"page,omitempty"`
	PerPage    string `json:"perPage,omitempty"`
	TotalPages string `json:"totalPages,omitempty"`
	Total      string `json:"total,omitempty"`
}
type Artists struct {
	Artist []Artist `json:"artist,omitempty"`
	Attr   Attr     `json:"@attr,omitempty"`
}

type Recenttracks struct {
	Attr  *RecenttracksAttr `json:"@attr,omitempty"`
	Track []Track           `json:"track,omitempty"`
}

type RecenttracksAttr struct {
	Page       string `json:"page,omitempty"`
	PerPage    string `json:"perPage,omitempty"`
	User       string `json:"user,omitempty"`
	Total      string `json:"total,omitempty"`
	TotalPages string `json:"totalPages,omitempty"`
}

type Track struct {
	Artist     *Album     `json:"artist,omitempty"`
	Attr       *TrackAttr `json:"@attr,omitempty"`
	Mbid       string     `json:"mbid,omitempty"`
	Album      *Album     `json:"album,omitempty"`
	Streamable string     `json:"streamable,omitempty"`
	URL        string     `json:"url,omitempty"`
	Name       string     `json:"name,omitempty"`
	Image      []Image    `json:"image,omitempty"`
	Date       *Date      `json:"date,omitempty"`
}

type Album struct {
	Mbid string `json:"mbid,omitempty"`
	Text string `json:"#text,omitempty"`
}

type TrackAttr struct {
	Nowplaying string `json:"nowplaying,omitempty"`
}

type Date struct {
	Uts  string `json:"uts,omitempty"`
	Text string `json:"#text,omitempty"`
}

type Image struct {
	Size *Size  `json:"size,omitempty"`
	Text string `json:"#text,omitempty"`
}

type Size string

type LastFMUser struct {
	User User `json:"user"`
}

type Registered struct {
	Unixtime string `json:"unixtime"`
	Text     int    `json:"#text"`
}
type User struct {
	Playlists  string     `json:"playlists"`
	Playcount  string     `json:"playcount"`
	Gender     string     `json:"gender"`
	Name       string     `json:"name"`
	Subscriber string     `json:"subscriber"`
	URL        string     `json:"url"`
	Country    string     `json:"country"`
	Image      []Image    `json:"image"`
	Registered Registered `json:"registered"`
	Type       string     `json:"type"`
	Age        string     `json:"age"`
	Bootstrap  string     `json:"bootstrap"`
	Realname   string     `json:"realname"`
}
