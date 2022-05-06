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
	GetPostByUser(userId int, start, end, count int) ([]*model.Post, error)
	// GetFollowerNum get number of followers
	GetFollowerNum(userId int) (int, error)
	// GetFolloweeNum get number of followee
	GetFolloweeNum(userId int) (int, error)
	// GetLikeNum get number of likes
	GetLikeNum(postId int) (int, error)
	// GetCommentNum get number of comments
	GetCommentNum(postId int) (int, error)
	// GetLikeUserByPost get who like the post
	GetLikeUserByPost(postId int) ([]*model.User, error)
	// GetCommentByPost get who comment the post
	GetCommentByPost(postId int) ([]*model.Comment, error)
	// GetCommentByID get comment detail by id
	GetCommentByID(commentId int) (*model.Comment, error)
	// GetPostLikeByUser post liked by user
	GetPostLikeByUser(userId int) ([]*model.Post, error)
	// GetPostStarByUser post  stared by user
	GetPostStarByUser(userId int) ([]*model.Post, error)
	// IsUserLikePost is user like the post
	IsUserLikePost(userId, postId int) (bool, error)
	// IsUserStarPost is user star the post
	IsUserStarPost(userId, postId int) (bool, error)
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
	// DelLike dislike the post
	DelLike(userId, postId int) error
	// DelComment delete the comment
	DelComment(commentId, postId int) error
	// DelPost delete the post
	DelPost(userId, postId int) error
}
