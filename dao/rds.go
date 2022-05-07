package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/heroyan/twitter/model"
	"strconv"
)

// redis 实现

type RedisDao struct {
	rdb *redis.Client
}

var ctx = context.Background()

func NewRedisDao(addr, passwd string, db int) *RedisDao {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	return &RedisDao{rdb: client}
}

// GetUserByID get user by id
func (rd *RedisDao) GetUserByID(userId int) (user *model.User, err error) {
	// uin:id hash saves the user info
	rt := model.User{}
	key := model.UinPrefix + fmt.Sprintf("%d", userId)
	ret := rd.rdb.HGetAll(ctx, key)
	err = ret.Scan(&rt)
	if err != nil {
		return nil, err
	}

	if rt.Id == 0 {
		return nil, nil
	}

	return &rt, nil
}

// GetUser get user by username
func (rd *RedisDao) GetUser(userName string) (user *model.User, err error) {
	key := model.UserPrefix + userName
	result, err := rd.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// user:username saves id in the redis
	// uin:id hash saves the user info
	rt := model.User{}
	key = model.UinPrefix + result
	ret := rd.rdb.HGetAll(ctx, key)
	err = ret.Scan(&rt)
	if err != nil {
		return nil, err
	}

	if rt.Id == 0 {
		return nil, nil
	}

	return &rt, nil
}

func (rd *RedisDao) generateId(modelType string) (int, error) {
	key := model.GenPrefix + modelType
	id, err := rd.rdb.Incr(ctx, key).Result()

	return int(id), err
}

func (rd *RedisDao) AddUser(user *model.User) (err error) {
	user.Id, err = rd.generateId(model.UserModel)
	if err != nil {
		return err
	}

	key := model.UserPrefix + user.UserName
	_, err = rd.rdb.Set(ctx, key, user.Id, 0).Result()
	if err != nil {
		return err
	}
	key = model.UinPrefix + fmt.Sprintf("%d", user.Id)
	_, err = rd.rdb.HMSet(ctx, key, "id", user.Id, "name", user.Name, "user_name", user.UserName,
		"passwd", user.Passwd, "last_login_time", user.LastLoginTime, "nick", user.Nick, "session_id", user.SessionId,
		"age", user.Age, "create_time", user.CreateTime).Result()

	return
}

func (rd *RedisDao) GetPost(id int) (post *model.Post, err error) {
	key := model.PostPrefix + fmt.Sprintf("%d", id)
	cmd := rd.rdb.HGetAll(ctx, key)
	if err != nil {
		return nil, err
	}
	rt := model.Post{}
	err = cmd.Scan(&rt)
	if err != nil {
		return nil, err
	}
	if rt.Id == 0 {
		return nil, nil
	}

	return &rt, nil
}

func (rd *RedisDao) AddPost(post *model.Post) (err error) {
	post.Id, err = rd.generateId(model.PostModel)
	if err != nil {
		return err
	}

	key := model.PostPrefix + fmt.Sprintf("%d", post.Id)
	_, err = rd.rdb.HMSet(ctx, key, "id", post.Id, "title", post.Title, "content", post.Content,
		"user_id", post.UserId, "create_time", post.CreateTime).Result()
	if err != nil {
		return err
	}

	// add to the user's post list
	// the newest one at the head
	key = model.UserPostPrefix + fmt.Sprintf("%d", post.UserId)
	_, err = rd.rdb.LPush(ctx, key, post.Id).Result()

	return
}

func (rd *RedisDao) GetPostByUser(userId, start, count int) (postList []*model.Post, err error) {
	key := model.UserPostPrefix + fmt.Sprintf("%d", userId)
	postIds, err := rd.rdb.LRange(ctx, key, int64(start), int64(start+count)).Result()
	if err != nil {
		return nil, err
	}
	for _, postId := range postIds {
		id, err := strconv.Atoi(postId)
		if err != nil {
			return nil, err
		}
		post, err := rd.GetPost(id)
		if err != nil {
			return nil, err
		}
		postList = append(postList, post)
	}

	return
}

// GetFollowerNum get how many users follow me
func (rd *RedisDao) GetFollowerNum(userId int) (int, error) {
	key := model.UserFollowerPrefix + fmt.Sprintf("%d", userId)
	num, err := rd.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// GetFolloweeNum get how many users I follow
func (rd *RedisDao) GetFolloweeNum(userId int) (int, error) {
	key := model.UserFolloweePrefix + fmt.Sprintf("%d", userId)
	num, err := rd.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// GetLikeNum get number of likes
func (rd *RedisDao) GetLikeNum(postId int) (int, error) {
	key := model.PostLikePrefix + fmt.Sprintf("%d", postId)
	num, err := rd.rdb.SCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// GetCommentNum get number of comments
func (rd *RedisDao) GetCommentNum(postId int) (int, error) {
	key := model.PostCommentPrefix + fmt.Sprintf("%d", postId)
	num, err := rd.rdb.LLen(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	return int(num), nil
}

// GetLikeUserByPost one user only can like one post once
// a post likes may be a large amount, here just give random 100
func (rd *RedisDao) GetLikeUserByPost(postId int) ([]*model.User, error) {
	var maxRandCount int64 = 100
	key := model.PostLikePrefix + fmt.Sprintf("%d", postId)
	uidList, err := rd.rdb.SRandMemberN(ctx, key, maxRandCount).Result()
	if err != nil {
		return nil, err
	}
	var userList []*model.User
	for _, uid := range uidList {
		userId, err := strconv.Atoi(uid)
		if err != nil {
			return nil, err
		}
		user, err := rd.GetUserByID(userId)
		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	return userList, nil
}

// GetCommentByPost get comments by the post
func (rd *RedisDao) GetCommentByPost(postId, start, count int) ([]*model.Comment, error) {
	key := model.PostCommentPrefix + fmt.Sprintf("%d", postId)
	commentIds, err := rd.rdb.LRange(ctx, key, int64(start), int64(start+count)).Result()
	if err != nil {
		return nil, err
	}

	var commentList []*model.Comment
	for _, commentId := range commentIds {
		key = model.CommentDetailPrefix + commentId
		var cmt model.Comment
		cmd := rd.rdb.HGetAll(ctx, key)
		err = cmd.Scan(&cmt)
		if err != nil {
			return nil, err
		}
		commentList = append(commentList, &cmt)
	}

	return commentList, nil
}

// AddLike someone likes the post
func (rd *RedisDao) AddLike(like *model.Like) error {
	key := model.PostLikePrefix + fmt.Sprintf("%d", like.PostId)
	_, err := rd.rdb.SAdd(ctx, key, like.UserId).Result()

	// add to someone's like list
	key = model.UserLikePrefix + fmt.Sprintf("%d", like.UserId)
	_, err = rd.rdb.ZAdd(ctx, key, &redis.Z{Score: float64(like.CreateTime), Member: like.PostId}).Result()

	return err
}

// AddComment someone comments the post
func (rd *RedisDao) AddComment(comment *model.Comment) (err error) {
	comment.Id, err = rd.generateId(model.CommentModel)
	if err != nil {
		return err
	}

	// first add detail, then add id map
	// the best way is to use reids pipeline
	pipe := rd.rdb.Pipeline()
	key := model.CommentDetailPrefix + fmt.Sprintf("%d", comment.Id)
	pipe.HMSet(ctx, key, "id", comment.Id, "post_id", comment.PostId, "user_id", comment.UserId,
		"content", comment.Content, "create_time", comment.CreateTime)

	key = model.PostCommentPrefix + fmt.Sprintf("%d", comment.PostId)
	pipe.LPush(ctx, key, comment.Id)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

// AddStar someone star the post
func (rd *RedisDao) AddStar(star *model.Star) error {
	pipe := rd.rdb.Pipeline()
	key := model.PostStarPrefix + fmt.Sprintf("%d", star.PostId)
	pipe.SAdd(ctx, key, star.UserId)

	key = model.UserStarPrefix + fmt.Sprintf("%d", star.UserId)
	// add to someone's star list
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(star.CreateTime), Member: star.PostId})
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

// AddFollower someone follows the other
func (rd *RedisDao) AddFollower(follow *model.Follow) error {
	// add to my fans list
	key := model.UserFollowerPrefix + fmt.Sprintf("%d", follow.FolloweeId)
	_, err := rd.rdb.SAdd(ctx, key, follow.FollowerId).Result()
	if err != nil {
		return err
	}
	// add to my followee list
	key = model.UserFolloweePrefix + fmt.Sprintf("%d", follow.FollowerId)
	_, err = rd.rdb.SAdd(ctx, key, follow.FolloweeId).Result()

	return err
}

// GetCommentByID get comment detail by id
func (rd *RedisDao) GetCommentByID(commentId int) (*model.Comment, error) {
	key := model.CommentDetailPrefix + fmt.Sprintf("%d", commentId)
	cmd := rd.rdb.HGetAll(ctx, key)
	var comment model.Comment
	err := cmd.Scan(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (rd *RedisDao) getPostByModel(modelType string, userId, start, end, count int) ([]*model.Post, error) {
	key := modelType + fmt.Sprintf("%d", userId)
	postIds, err := rd.rdb.ZRangeByScore(ctx, key,
		&redis.ZRangeBy{
			Min:    fmt.Sprintf("%d", start),
			Max:    fmt.Sprintf("%d", end),
			Count:  int64(count),
			Offset: 0,
		}).Result()
	if err != nil {
		return nil, err
	}

	var postList []*model.Post
	for _, postId := range postIds {
		id, err := strconv.Atoi(postId)
		if err != nil {
			return nil, err
		}
		post, err := rd.GetPost(id)
		if err != nil {
			return nil, err
		}
		postList = append(postList, post)
	}

	return postList, nil
}

// GetPostLikeByUser post liked by user
func (rd *RedisDao) GetPostLikeByUser(userId, start, end, count int) ([]*model.Post, error) {
	return rd.getPostByModel(model.UserLikePrefix, userId, start, end, count)
}

// GetPostStarByUser post  stared by user
func (rd *RedisDao) GetPostStarByUser(userId, start, end, count int) ([]*model.Post, error) {
	return rd.getPostByModel(model.UserStarPrefix, userId, start, end, count)
}

func (rd *RedisDao) isModelPost(modelType string, userId, postId int) (bool, error) {
	key := modelType + fmt.Sprintf("%d", userId)
	_, err := rd.rdb.ZScore(ctx, key, fmt.Sprintf("%d", postId)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

// IsUserLikePost is user like the post
func (rd *RedisDao) IsUserLikePost(userId, postId int) (bool, error) {
	return rd.isModelPost(model.UserLikePrefix, userId, postId)
}

// IsUserStarPost is user star the post
func (rd *RedisDao) IsUserStarPost(userId, postId int) (bool, error) {
	return rd.isModelPost(model.UserStarPrefix, userId, postId)
}

// IsUserNameExists is username exists
func (rd *RedisDao) IsUserNameExists(userName string) (bool, error) {
	key := model.UserPrefix + userName
	_, err := rd.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return true, err
	}

	return true, nil
}

// DelLike dislike the post
func (rd *RedisDao) DelLike(userId, postId int) error {
	key := model.PostLikePrefix + fmt.Sprintf("%d", postId)
	_, err := rd.rdb.SRem(ctx, key, userId).Result()
	if err != nil {
		return err
	}

	key = model.UserLikePrefix + fmt.Sprintf("%d", userId)
	_, err = rd.rdb.ZRem(ctx, key, postId).Result()

	return err
}

// DelComment delete the comment
func (rd *RedisDao) DelComment(commentId, postId int) error {
	key := model.PostCommentPrefix + fmt.Sprintf("%d", postId)
	_, err := rd.rdb.LRem(ctx, key, 1, commentId).Result()
	if err != nil {
		return err
	}

	key = model.CommentDetailPrefix + fmt.Sprintf("%d", commentId)
	_, err = rd.rdb.Del(ctx, key).Result()

	return err
}

// DelPost delete the post
func (rd *RedisDao) DelPost(userId, postId int) error {
	key := model.PostPrefix + fmt.Sprintf("%d", postId)
	_, err := rd.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	key = model.UserPostPrefix + fmt.Sprintf("%d", userId)
	_, err = rd.rdb.LRem(ctx, key, 1, postId).Result()

	return err
}
