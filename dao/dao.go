package dao

import "github.com/heroyan/twitter/model"

type Dao interface {
	// GetUser get user info
	GetUser(userName string) (user *model.User, err error)
	// GetUserByID get user info
	GetUserByID(userId int) (user *model.User, err error)
	// GetPost get detail of the post
	GetPost(id int) (post *model.Post, err error)
	// GetPostByUser get user's post list
	GetPostByUser(userId int, start, count int) ([]*model.Post, error)
	// GetFollowerNum get number of followers
	GetFollowerNum(userId int) (int, error)
	// GetFolloweeNum get number of followee
	GetFolloweeNum(userId int) (int, error)
	// GetStarNum get number of stars
	GetStarNum(postId int) (int, error)
	// GetLikeNum get number of likes
	GetLikeNum(postId int) (int, error)
	// GetCommentNum get number of comments
	GetCommentNum(postId int) (int, error)
	// GetLikeUserByPost get who like the post
	GetLikeUserByPost(postId int) ([]*model.User, error)
	// GetCommentByPost get who comment the post
	GetCommentByPost(postId, start, count int) ([]*model.Comment, error)
	// GetCommentByID get comment detail by id
	GetCommentByID(commentId int) (*model.Comment, error)
	// GetPostLikeByUser post liked by user
	GetPostLikeByUser(userId, start, count int) ([]*model.Post, error)
	// GetPostStarByUser post  stared by user
	GetPostStarByUser(userId, start, count int) ([]*model.Post, error)
	// GetPostFollowByUser post followed by user
	GetPostFollowByUser(userId, start, count int) ([]*model.Post, error)
	// GetHotPost hot post recommende to user
	GetHotPost(userId, count int) ([]*model.Post, error)
	// GetFollowers who follow me
	GetFollowers(userId, count int) ([]*model.User, error)
	// GetFollowees I follow who
	GetFollowees(userId, count int) ([]*model.User, error)
	// IsUserLikePost is user like the post
	IsUserLikePost(userId, postId int) (bool, error)
	// IsUserStarPost is user star the post
	IsUserStarPost(userId, postId int) (bool, error)
	// IsUserFollow is user follow the other
	IsUserFollow(followerId, followeeId int) (bool, error)
	// IsUserNameExists is username exists
	IsUserNameExists(userName string) (bool, error)
	// AddPost save post to storage
	AddPost(post *model.Post) (err error)
	// AddUser save user info
	AddUser(user *model.User) (err error)
	// AddLike someone likes the post
	AddLike(like *model.Like) error
	// AddComment someone comments the post
	AddComment(comment *model.Comment) error
	// AddStar someone star the post
	AddStar(star *model.Star) error
	// AddFollower someone follows the other
	AddFollower(follow *model.Follow) error
	// UnFollow someone
	UnFollow(follow *model.Follow) error
	// DelLike dislike the post
	DelLike(userId, postId int) error
	// DelStar un star the post
	DelStar(userId, postId int) error
	// DelComment delete the comment
	DelComment(commentId, postId int) error
	// DelPost delete the post
	DelPost(userId, postId int) error
	// GetSessionUser session resolve
	GetSessionUser(sessionId string) (*model.User, error)
	// SetSessionUser set session user
	SetSessionUser(sessionId string, userId int, expire int) error
	// DelSession delete the session
	DelSession(sessionId string) error
}
