package service

import (
	"github.com/heroyan/twitter/dao"
	"github.com/heroyan/twitter/model"
)

type UserService struct {
	daoObj *dao.Dao
}

func RegisterUser(user *model.User) (id int, err error) {
	return id, err
}

func LoginUser(user *model.User) (err error) {
	return err
}
