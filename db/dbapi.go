package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var Conn *gorm.DB

func InitConnection(dbfile string) error {
	Conn, err := gorm.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}
	defer Conn.Close()

	//TODO Build Debug Model
	Conn.LogMode(true)
	Conn.AutoMigrate(
		User{},
		OnlineList{},
	)

	return nil
}

func GetUserByName(username string) (User, error) {
	user := User{}
	err := Conn.First(&user, "username = ?", username).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetOnlineList() ([]OnlineList, error) {
	list := &[]OnlineList{}

	err := Conn.Find(list).Error
	if err != nil {
		return nil, err
	}

	return *list, nil
}

func CreateUser(user *User) error {
	err := Conn.Create(user).Error
	if err != nil {
		return err
	}

	return nil
}

func AddUser2List(ol *OnlineList) error {
	err := Conn.Create(ol).Error
	if err != nil {
		return err
	}

	return nil
}

func KickOutUser(mac string) error {
	err := Conn.Delete(User{}, "mac = ?", mac).Error
	if err != nil {
		return err
	}
	return nil
}
