package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
	"log"
)

var Conn *gorm.DB

func InitConnection(dbfile string) error {
	var err error
	Conn, err = gorm.Open("sqlite3", dbfile)
	if err != nil {
		return err
	}

	//TODO Build Debug Model
	Conn.LogMode(true)
	Conn.AutoMigrate(
		User{},
		OnlineUser{},
	)
	return nil
}

func GetUserByName(username string) (User, error) {
	user := User{}
	err := Conn.Find(&user, "username = ?", username).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetOnlineUserList() ([]OnlineUser, error) {
	list := &[]OnlineUser{}

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

func AddUser2List(ol *OnlineUser) error {
	/*
	 * make sure there is only one login record
	 */
	Conn.Delete(ol, "username = ?", ol.Username)
	err := Conn.Create(ol).Error
	if err != nil {
		return err
	}

	return nil
}

func KickOutUser(username, mac string) error {
	err := Conn.Delete(OnlineUser{}, "username = ? AND mac = ?", username, mac).Error
	if err != nil {
		return err
	}
	return nil
}

func CleanExpiredUserList() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			err := Conn.Delete(&OnlineUser{}, "expired_time_stamp - ? <= 0", time.Now().Unix()).Error
			if err != nil {
				log.Println("本次定时清理任务执行失败")
			} else {
				log.Println("定时清理任务执行成功")
			}
		}
	}
}
