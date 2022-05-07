package service

import (
	"time"

	"github.com/heroyan/twitter/dao"
	"github.com/heroyan/twitter/model"
)

type PostService struct {
	daoObj dao.Dao
}

func NewPostService(obj dao.Dao) *PostService {
	return &PostService{daoObj: obj}
}

func (svc *PostService) AddPost(post *model.Post) (int, error) {
	post.CreateTime = int(time.Now().Unix())
	err := svc.daoObj.AddPost(post)

	return post.Id, err
}

func (svc *PostService) DelPost(userId, postId int) error {
	return svc.daoObj.DelPost(userId, postId)
}

func (svc *PostService) GetPost(postId int) (*model.Post, error) {
	return svc.daoObj.GetPost(postId)
}

func (svc *PostService) AddComment(comment *model.Comment) error {
	comment.CreateTime = int(time.Now().Unix())
	err := svc.daoObj.AddComment(comment)

	return err
}

func (svc *PostService) DelComment(commentId, postId int) error {
	return svc.daoObj.DelComment(commentId, postId)
}

func (svc *PostService) GetCommentByPost(postId, start, count int) ([]*model.Comment, error) {
	return svc.daoObj.GetCommentByPost(postId, start, count)
}
