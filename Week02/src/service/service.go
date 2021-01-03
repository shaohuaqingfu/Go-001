package service

import (
	"Week02/src/dao"
	"Week02/src/model"
)

type GeekErrService struct {
	dao dao.UserDao
}

func (svc *GeekErrService) GetById(id string) (*model.User, error) {
	user, err := svc.dao.GetById(id)
	if err != nil {
		if svc.dao.Exists(err) {
			//doSomeThings()
		}
		return nil, err
	}
	return user, nil
}
