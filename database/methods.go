package database

func (u *User) IsValid() bool {
	return u != nil && u.UserID != 0 && u.LastFmUsername != ""
}
