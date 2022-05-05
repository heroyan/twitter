package dao

import "github.com/heroyan/twitter/model"

type Dao interface {
	GetUser(userName string) (user *model.User, err error)
	SaveUser(user *model.User) (err error)
	GetPost(id int) (post *model.Post, err error)
	SavePost(post *model.Post) (err error)
}
