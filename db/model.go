package db

import "time"

type BaseModel struct {
	Id        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	BaseModel
	Username string `gorm:"not null;unique" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Level    int    `gorm:"default:'1'" json:"level"`
}

type OnlineList struct {
	BaseModel        `json:"-"`
	Username         string    `json:"username"`
	IP               string    `json:"ip"`
	Mac              string    `json:"mac"`
	ExpiredAt        time.Time `json:"expired_at"`
	ExpiredTimeStamp int64     `json:"-"`
}
