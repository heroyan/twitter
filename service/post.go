package service

import (
	"errors"
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

func (svc *PostService) AddPost(post *model.Post) error {
	post.CreateTime = int(time.Now().Unix())
	err := svc.daoObj.AddPost(post)

	return err
}

func (svc *PostService) DelPost(userId, postId int) error {
	// must check if userId matches with the post's userId
	post, err := svc.GetPost(postId)
	if err != nil {
		return err
	}
	if post == nil || post.UserId != userId {
		return errors.New("post not found")
	}

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

func (svc *PostService) DelComment(commentId, userId int) error {
	// check if the comment is published by the user
	cmt, err := svc.daoObj.GetCommentByID(commentId)
	if err != nil {
		return err
	}
	if cmt == nil || cmt.UserId != userId {
		return errors.New("comment not found")
	}

	return svc.daoObj.DelComment(commentId, cmt.PostId)
}

func (svc *PostService) GetCommentByPost(postId, start, count int) ([]*model.Comment, error) {
	return svc.daoObj.GetCommentByPost(postId, start, count)
}

func (svc *PostService) DelLike(userId, postId int) error {
	return svc.daoObj.DelLike(userId, postId)
}

func (svc *PostService) DelStar(userId, postId int) error {
	return svc.daoObj.DelStar(userId, postId)
}

func (svc *PostService) AddLike(like *model.Like) error {
	yes, err := svc.daoObj.IsUserLikePost(like.UserId, like.PostId)
	if err != nil {
		return err
	}

	// if already liked, then skip
	if yes {
		return nil
	}

	return svc.daoObj.AddLike(like)
}

func (svc *PostService) AddStar(star *model.Star) error {
	yes, err := svc.daoObj.IsUserStarPost(star.UserId, star.PostId)
	if err != nil {
		return err
	}

	// if already liked, then skip
	if yes {
		return nil
	}

	return svc.daoObj.AddStar(star)
}

func (svc *PostService) GetLikeNum(postId int) (int, error) {
	return svc.daoObj.GetLikeNum(postId)
}

func (svc *PostService) GetCommentNum(postId int) (int, error) {
	return svc.daoObj.GetCommentNum(postId)
}

func (svc *PostService) GetStarNum(postId int) (int, error) {
	return svc.daoObj.GetStarNum(postId)
}

func (svc *PostService) IsUserStarPost(userId, postId int) (bool, error) {
	return svc.daoObj.IsUserStarPost(userId, postId)
}

func (svc *PostService) IsUserLikePost(userId, postId int) (bool, error) {
	return svc.daoObj.IsUserLikePost(userId, postId)
}
