package auth

type SpotifyUser struct {
	UserId       int64  `gorm:"primaryKey"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
