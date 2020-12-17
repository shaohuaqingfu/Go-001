package model

import "time"

type Order struct {
	Id           int32     `gorm:"id;primary_key"`
	OrderNo      string    `gorm:"order_no"`
	CreateTime   time.Time `gorm:"create_time"`
	CreateUserId int32     `gorm:"create_user_id"`
}
