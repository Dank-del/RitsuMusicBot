package odesli

type OdesliResponse struct {
	EntityUniqueID     string                        `json:"entityUniqueId,omitempty"`
	UserCountry        string                        `json:"userCountry,omitempty"`
	PageURL            string                        `json:"pageUrl,omitempty"`
	EntitiesByUniqueID map[string]EntitiesByUniqueID `json:"entitiesByUniqueId,omitempty"`
	LinksByPlatform    *LinksByPlatform              `json:"linksByPlatform,omitempty"`
	StatusCode         int64                         `json:"statusCode,omitempty"`
	Code               string                        `json:"code,omitempty"`
}

type EntitiesByUniqueID struct {
	ID              string   `json:"id,omitempty"`
	Type            string   `json:"type,omitempty"`
	Title           string   `json:"title,omitempty"`
	ArtistName      string   `json:"artistName,omitempty"`
	ThumbnailURL    string   `json:"thumbnailUrl,omitempty"`
	ThumbnailWidth  int64    `json:"thumbnailWidth,omitempty"`
	ThumbnailHeight int64    `json:"thumbnailHeight,omitempty"`
	APIProvider     string   `json:"apiProvider,omitempty"`
	Platforms       []string `json:"platforms,omitempty"`
}

type LinksByPlatform struct {
	AppleMusic   *AppleMusic    `json:"appleMusic,omitempty"`
	Itunes       *AppleMusic    `json:"itunes,omitempty"`
	Napster      *OtherPlatform `json:"napster,omitempty"`
	Soundcloud   *OtherPlatform `json:"soundcloud,omitempty"`
	Spotify      *OtherPlatform `json:"spotify,omitempty"`
	Youtube      *OtherPlatform `json:"youtube,omitempty"`
	YoutubeMusic *OtherPlatform `json:"youtubeMusic,omitempty"`
	AmazonMusic  *OtherPlatform `json:"amazonMusic,omitempty"`
	AmazonStore  *OtherPlatform `json:"amazonStore,omitempty"`
	Deezer       *OtherPlatform `json:"deezer,omitempty"`
	Pandora      *OtherPlatform `json:"pandora,omitempty"`
	Tidal        *OtherPlatform `json:"tidal,omitempty"`
	Yandex       *OtherPlatform `json:"yandex,omitempty"`
}

type AppleMusic struct {
	Country             string `json:"country,omitempty"`
	URL                 string `json:"url,omitempty"`
	NativeAppURIMobile  string `json:"nativeAppUriMobile,omitempty"`
	NativeAppURIDesktop string `json:"nativeAppUriDesktop,omitempty"`
	EntityUniqueID      string `json:"entityUniqueId,omitempty"`
}

type OtherPlatform struct {
	Country             string `json:"country,omitempty"`
	URL                 string `json:"url,omitempty"`
	EntityUniqueID      string `json:"entityUniqueId,omitempty"`
	NativeAppURIDesktop string `json:"nativeAppUriDesktop,omitempty"`
}
