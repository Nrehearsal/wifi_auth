package db

import "time"

type BaseModel struct {
	Id        int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type User struct {
	BaseModel
	Usernmae       string `gorm:"not null;unique" json:"usernmae"`
	Password       string `gorm:"not null" json:"-"`
	AccessDuration int    `gorm:"default:'720'" json:"access_duration"`
}

type OnlineList struct {
	BaseModel
	UserId int    `json:"user_id"`
	IP     string `json:"ip"`
	Mac    string `json:"mac"`
	//per hour
	Lifetime int `json:"lifetime"`
}
