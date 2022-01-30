package last_fm

const (
	Extralarge Size = "extralarge"
	// Large      Size = "large"
	// Medium     Size = "medium"
	// Small      Size = "small"
)

const (
	recentTracksBaseUrl        = "https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&user="
	userBaseUrl                = "https://ws.audioscrobbler.com/2.0/?method=user.getinfo&user="
	getTopArtistsBaseURl       = "https://ws.audioscrobbler.com/2.0/?method=chart.gettopartists"
	getTrackInfoBaseUrl        = "https://ws.audioscrobbler.com/2.0/?method=track.getInfo"
	getTopTracksForUserBaseUrl = "https://ws.audioscrobbler.com/2.0/?method=user.getTopTracks&format=json&period=overall&user="
)
