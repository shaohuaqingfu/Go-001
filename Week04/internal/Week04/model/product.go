package model

import "time"

type Product struct {
	Id           int32     `gorm:"id;primary_key"`
	Name         string    `gorm:"name"`
	CreateTime   time.Time `gorm:"create_time"`
	CreateUserId int32     `gorm:"create_user_id"`
}
