package last_fm

type GetRecentTracks struct {
	Recenttracks *Recenttracks `json:"recenttracks,omitempty"`
	Error        int           `json:"error,omitempty"`
	Message      string        `json:"message,omitempty"`
}

type TopTrackResponse struct {
	Toptracks *Toptracks `json:"toptracks,omitempty"`
	Error     int        `json:"error,omitempty"`
	Message   string     `json:"message,omitempty"`
}

type Toptracks struct {
	Attr  *ToptracksAttr `json:"@attr,omitempty"`
	Track []Track        `json:"track,omitempty"`
}

type ToptracksAttr struct {
	Page       string `json:"page,omitempty"`
	PerPage    string `json:"perPage,omitempty"`
	User       string `json:"user,omitempty"`
	Total      string `json:"total,omitempty"`
	TotalPages string `json:"totalPages,omitempty"`
}

type GetLastFmTrackResponse struct {
	Track   TrackInfo `json:"track"`
	Error   int       `json:"error,omitempty"`
	Message string    `json:"message,omitempty"`
}

type Streamable struct {
	Text      string `json:"#text,omitempty"`
	Fulltrack string `json:"fulltrack,omitempty"`
}

type GetTopArtistsResponse struct {
	Artists Artists `json:"artists,omitempty"`
	Error   int     `json:"error,omitempty"`
	Message string  `json:"message,omitempty"`
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
	Position   string `json:"position"`
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

type TrackInfo struct {
	Name          string      `json:"name"`
	Mbid          string      `json:"mbid"`
	URL           string      `json:"url"`
	Duration      string      `json:"duration"`
	Streamable    *Streamable `json:"streamable"`
	Listeners     string      `json:"listeners"`
	Playcount     string      `json:"playcount"`
	Artist        *Artist     `json:"artist"`
	Album         *Album      `json:"album"`
	Toptags       *Toptags    `json:"toptags"`
	Wiki          *Wiki       `json:"wiki"`
	Userplaycount string      `json:"userplaycount,omitempty"`
}

type Track struct {
	Attr   *TrackAttr `json:"@attr,omitempty"`
	Artist *Artist    `json:"artist,omitempty"`
	Mbid   string     `json:"mbid,omitempty"`
	Image  []Image    `json:"image,omitempty"`
	URL    string     `json:"url,omitempty"`
	// Streamable Streamable `json:"streamable,omitempty"`
	Listeners string   `json:"listeners,omitempty"`
	Playcount string   `json:"playcount,omitempty"`
	Album     *Album   `json:"album,omitempty"`
	Name      string   `json:"name,omitempty"`
	Loved     string   `json:"loved,omitempty"`
	Date      *Date    `json:"date,omitempty"`
	Toptags   *Toptags `json:"toptags,omitempty"`
	Wiki      *Wiki    `json:"wiki,omitempty"`
	Duration  string   `json:"duration,omitempty"`
}

type Album struct {
	Mbid   string  `json:"mbid,omitempty"`
	Text   string  `json:"#text,omitempty"`
	Artist string  `json:"artist,omitempty"`
	Title  string  `json:"title,omitempty"`
	URL    string  `json:"url,omitempty"`
	Image  []Image `json:"image,omitempty"`
	Attr   Attr    `json:"@attr,omitempty"`
}

type TrackAttr struct {
	Nowplaying string `json:"nowplaying,omitempty"`
}

type Date struct {
	Uts  string `json:"uts,omitempty"`
	Text string `json:"#text,omitempty"`
}

type Image struct {
	Size Size   `json:"size,omitempty"`
	Text string `json:"#text,omitempty"`
}

type Size string

type LastFMUser struct {
	User *User `json:"user"`
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

type Toptags struct {
	Tag []Tag `json:"tag"`
}
type Wiki struct {
	Published string `json:"published"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
}

type Tag struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
