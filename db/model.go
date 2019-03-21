package db

import "time"

type BaseModel struct {
	Id        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type User struct {
	BaseModel
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Level    int    `gorm:"default:'1'" json:"level"`
}

type OnlineList struct {
	BaseModel
	UserId int    `json:"user_id"`
	IP     string `json:"ip"`
	Mac    string `json:"mac"`
	//per hour
	ExpiredAt        time.Time `json:"expired_at"`
	ExpiredTimeStamp int64     `json:"expired_time_stamp"`
}
