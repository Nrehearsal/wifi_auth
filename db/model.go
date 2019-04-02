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

type OnlineUser struct {
	BaseModel        `json:"-"`
	Username         string    `json:"username" gorm:"not null;unique"`
	Level            int       `json:"level" gorm:"not null"`
	IP               string    `json:"ip" gorm:"not null"`
	Mac              string    `json:"mac" gorm:"not null"`
	ExpiredAt        time.Time `json:"expired_at" gorm:"not null"`
	ExpiredTimeStamp int64     `json:"-" gorm:"not null"`
}
