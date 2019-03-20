package db

import (
	"github.com/jinzhu/gorm"
	"tech.feiyuapi.com/feiyu.ci/client-be/api/db"
)

var Conn *gorm.DB

func InitConnection(dbfile string) error {
	Conn, err := gorm.Open("sqlite3.db", dbfile)
	if err != nil {
		return err
	}
	defer Conn.Close()

	//TODO Build Debug Model
	Conn.LogMode(true)
	Conn.AutoMigrate()

	return nil
}

func GetUserByName(username string) (User, error) {
	user := User{}
	err := db.Conn.First(&user, "username = ?", username).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
