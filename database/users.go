package database

import (
	"errors"
	"strings"
)

type User struct {
	UserID int64 `gorm:"primaryKey" json:"user_id"`
	LastFmUsername string
}


func UpdateUser(UserID int64, LastFmUsername string) {
	tx := SESSION.Begin()
	user := &User{UserID: UserID, LastFmUsername: strings.ToLower(LastFmUsername)}
	tx.Save(user)
	tx.Commit()
}

func GetUser(UserID int64) (u *User, err error){
	if SESSION == nil {
		return nil, errors.New("cannot access to SESSION " +
			"of db, because it's nil")
	}

	p := User{}
	SESSION.Where("user_id = ?", UserID).Take(&p)
	return &p, nil
}

type Tiddie struct {
	Size int64
	Weight int64
}