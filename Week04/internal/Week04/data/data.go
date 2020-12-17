package data

import (
	"Week04/internal/Week04/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type OrderData struct {
	DB *gorm.DB
}

func (d *OrderData) GetById(id int32) (*model.Order, error) {
	order := &model.Order{}

	if err := d.DB.Table("t_order").Where("id = ?", id).Find(order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(gorm.ErrRecordNotFound, "没有找到订单")
		}
		return nil, errors.Wrap(err, "数据库查询失败")
	}
	return order, nil
}

func (d *OrderData) IsNotExists(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (d *OrderData) GetByCreateUserId(createUserId int32) ([]*model.Order, error) {
	var orders []*model.Order
	if err := d.DB.Table("t_order").Where("create_user_id = ?", createUserId).Find(&orders).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(gorm.ErrRecordNotFound, "没有找到订单")
		}
		return nil, errors.Wrap(err, "数据库查询失败")
	}
	return orders, nil
}
