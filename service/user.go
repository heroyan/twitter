package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/heroyan/twitter/dao"
	"github.com/heroyan/twitter/model"
)

type UserService struct {
	daoObj dao.Dao
}

func NewUserService(obj dao.Dao) *UserService {
	return &UserService{daoObj: obj}
}

func getMd5(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("%x", hash)
}

func (svc *UserService) RegisterUser(user *model.User) (id int, err error) {
	// if username exists
	exists, err := svc.daoObj.IsUserNameExists(user.UserName)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, errors.New("username already exists")
	}
	// passwd must be hashed to store
	user.Passwd = getMd5(user.Passwd)
	// auto add register time
	user.CreateTime = int(time.Now().Unix())
	err = svc.daoObj.AddUser(user)

	return user.Id, err
}

func (svc *UserService) LoginUser(user *model.User) (err error) {
	// if username and passwd matches
	u, err := svc.daoObj.GetUser(user.UserName)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.New("username or passwd wrong")
	}
	if u.Passwd != getMd5(user.Passwd) {
		return errors.New("username or passwd wrong")
	}

	// used to set session
	user.Id = u.Id

	return err
}

func (svc *UserService) Logout(sessionId string) (err error) {
	return svc.daoObj.DelSession(sessionId)
}

func (svc *UserService) GetSessionUser(sessionId string) (*model.User, error) {
	return svc.daoObj.GetSessionUser(sessionId)
}

// SetSessionUser set session user
func (svc *UserService) SetSessionUser(sessionId string, userId, expire int) error {
	return svc.daoObj.SetSessionUser(sessionId, userId, expire)
}

func (svc *UserService) GetPostByUser(userId, start, count int) ([]*model.Post, error) {
	return svc.daoObj.GetPostByUser(userId, start, count)
}

func (svc *UserService) GetPostLikeByUser(userId, start, count int) ([]*model.Post, error) {
	return svc.daoObj.GetPostLikeByUser(userId, start, count)
}

func (svc *UserService) GetPostStarByUser(userId, start, count int) ([]*model.Post, error) {
	return svc.daoObj.GetPostStarByUser(userId, start, count)
}

func (svc *UserService) GetPostFollowByUser(userId, start, count int) ([]*model.Post, error) {
	return svc.daoObj.GetPostFollowByUser(userId, start, count)
}

func (svc *UserService) GetHotPost(userId, count int) ([]*model.Post, error) {
	return svc.daoObj.GetHotPost(userId, count)
}
