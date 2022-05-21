package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/heroyan/twitter/model"
	"strconv"
	"time"
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
	if user.Id == 0 {
		user.Id, err = rd.generateId(model.UserModel)
		if err != nil {
			return err
		}
	}

	pipe := rd.rdb.Pipeline()
	key := model.UserPrefix + user.UserName
	pipe.Set(ctx, key, user.Id, 0)

	key = model.UinPrefix + fmt.Sprintf("%d", user.Id)
	pipe.HMSet(ctx, key, "id", user.Id, "name", user.Name, "user_name", user.UserName,
		"passwd", user.Passwd, "last_login_time", user.LastLoginTime, "nick", user.Nick, "session_id", user.SessionId,
		"age", user.Age, "create_time", user.CreateTime)

	//ignore results
	_, err = pipe.Exec(ctx)

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

	rt.LikeNum, _ = rd.GetLikeNum(rt.Id)
	rt.StarNum, _ = rd.GetStarNum(rt.Id)
	rt.CommentNum, _ = rd.GetCommentNum(rt.Id)
	// ignore the errors
	user, _ := rd.GetUserByID(rt.UserId)
	if user != nil {
		rt.UserNick = user.Nick
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

	// add to allpost model, ignore the error or add a log info
	rd.rdb.ZAdd(ctx, model.AllPostModel, &redis.Z{
		Score:  float64(post.CreateTime),
		Member: post.Id,
	})

	// add to the follower's timeline, this should be done background, we do this here for simple
	idList, err := rd.GetAllFollowers(post.UserId)
	if err != nil {
		return err
	}
	for _, id := range idList {
		key = model.UserTimelinePrefix + id
		// ignore errors or add a log info
		rd.rdb.LPush(ctx, key, post.Id)
	}

	return
}

// GetAllFollowers get who follow me
func (rd *RedisDao) GetAllFollowers(userId int) ([]string, error) {
	key := model.UserFollowerPrefix + fmt.Sprintf("%d", userId)
	// if there are a lot of followers, then it needs to pagination or get rand ones
	idList, err := rd.rdb.SMembers(ctx, key).Result()

	return idList, err
}

func (rd *RedisDao) getUserListByType(modelType string, userId, count int) ([]*model.User, error) {
	key := modelType + fmt.Sprintf("%d", userId)
	// if there are a lot of followers, then it needs to pagination or get rand ones
	idList, err := rd.rdb.SRandMemberN(ctx, key, int64(count)).Result()
	if err != nil {
		return nil, err
	}
	var userList []*model.User
	for _, id := range idList {
		// ignore errors
		uid, _ := strconv.Atoi(id)
		user, _ := rd.GetUserByID(uid)
		userList = append(userList, user)
	}

	return userList, nil
}

// GetFollowers get who follow me
func (rd *RedisDao) GetFollowers(userId, count int) ([]*model.User, error) {
	return rd.getUserListByType(model.UserFollowerPrefix, userId, count)
}

// GetFollowees get who follow me
func (rd *RedisDao) GetFollowees(userId, count int) ([]*model.User, error) {
	return rd.getUserListByType(model.UserFolloweePrefix, userId, count)
}

func (rd *RedisDao) GetPostByUser(userId, start, count int) (postList []*model.Post, err error) {
	key := model.UserPostPrefix + fmt.Sprintf("%d", userId)
	// here count must minus 1, because redis lrange contain the min and max pos
	postIds, err := rd.rdb.LRange(ctx, key, int64(start), int64(start+count-1)).Result()
	if err != nil {
		return nil, err
	}

	return rd.getPostByIdList(postIds)
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

// GetStarNum get number of likes
func (rd *RedisDao) GetStarNum(postId int) (int, error) {
	key := model.PostStarPrefix + fmt.Sprintf("%d", postId)
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
	// here count must minus 1, because redis lrange contain the min and max pos
	commentIds, err := rd.rdb.LRange(ctx, key, int64(start), int64(start+count-1)).Result()
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
		// get user nick name, ignore errors
		user, _ := rd.GetUserByID(cmt.UserId)
		if user != nil {
			cmt.UserNick = user.Nick
		}
		commentList = append(commentList, &cmt)
	}

	return commentList, nil
}

// AddLike someone likes the post
func (rd *RedisDao) AddLike(like *model.Like) error {
	pipe := rd.rdb.Pipeline()

	key := model.PostLikePrefix + fmt.Sprintf("%d", like.PostId)
	pipe.SAdd(ctx, key, like.UserId).Result()
	// add to someone's like list
	key = model.UserLikePrefix + fmt.Sprintf("%d", like.UserId)
	pipe.LPush(ctx, key, like.PostId)
	// set like flag to true
	key = model.IsLikePrefix + fmt.Sprintf("%d", like.UserId)
	pipe.SetBit(ctx, key, int64(like.PostId), 1)

	_, err := pipe.Exec(ctx)

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

	return err
}

// AddStar someone star the post
func (rd *RedisDao) AddStar(star *model.Star) error {
	pipe := rd.rdb.Pipeline()

	// add to the post star list
	key := model.PostStarPrefix + fmt.Sprintf("%d", star.PostId)
	pipe.SAdd(ctx, key, star.UserId)
	// add to someone's star list
	key = model.UserStarPrefix + fmt.Sprintf("%d", star.UserId)
	pipe.LPush(ctx, key, star.PostId)
	// set star flag to true
	key = model.IsStarPrefix + fmt.Sprintf("%d", star.UserId)
	pipe.SetBit(ctx, key, int64(star.PostId), 1)

	_, err := pipe.Exec(ctx)

	return err
}

// AddFollower someone follows the other
func (rd *RedisDao) AddFollower(follow *model.Follow) error {
	pipe := rd.rdb.Pipeline()
	// add to my fans list
	key := model.UserFollowerPrefix + fmt.Sprintf("%d", follow.FolloweeId)
	pipe.SAdd(ctx, key, follow.FollowerId)
	// add to my followee list
	key = model.UserFolloweePrefix + fmt.Sprintf("%d", follow.FollowerId)
	pipe.SAdd(ctx, key, follow.FolloweeId)
	// set follow flag
	key = model.IsFollowPrefix + fmt.Sprintf("%d", follow.FollowerId)
	pipe.SetBit(ctx, key, int64(follow.FolloweeId), 1)

	_, err := pipe.Exec(ctx)

	return err
}

// UnFollow unfollow someone
func (rd *RedisDao) UnFollow(follow *model.Follow) error {
	pipe := rd.rdb.Pipeline()
	// remove from my fans list
	key := model.UserFollowerPrefix + fmt.Sprintf("%d", follow.FolloweeId)
	pipe.SRem(ctx, key, follow.FollowerId)
	// remove from my followee list
	key = model.UserFolloweePrefix + fmt.Sprintf("%d", follow.FollowerId)
	pipe.SRem(ctx, key, follow.FolloweeId)
	// unset follow flag
	key = model.IsFollowPrefix + fmt.Sprintf("%d", follow.FollowerId)
	pipe.SetBit(ctx, key, int64(follow.FolloweeId), 0)

	_, err := pipe.Exec(ctx)

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

func (rd *RedisDao) getPostByModel(modelType string, userId, start, count int) ([]*model.Post, error) {
	key := modelType + fmt.Sprintf("%d", userId)
	postIds, err := rd.rdb.LRange(ctx, key, int64(start), int64(start+count-1)).Result()
	if err != nil {
		return nil, err
	}

	return rd.getPostByIdList(postIds)
}

// GetPostLikeByUser post liked by user
func (rd *RedisDao) GetPostLikeByUser(userId, start, count int) ([]*model.Post, error) {
	return rd.getPostByModel(model.UserLikePrefix, userId, start, count)
}

// GetPostStarByUser post  stared by user
func (rd *RedisDao) GetPostStarByUser(userId, start, count int) ([]*model.Post, error) {
	return rd.getPostByModel(model.UserStarPrefix, userId, start, count)
}

// GetPostFollowByUser post followed by user
func (rd *RedisDao) GetPostFollowByUser(userId, start, count int) ([]*model.Post, error) {
	return rd.getPostByModel(model.UserTimelinePrefix, userId, start, count)
}

func (rd *RedisDao) isModelPost(modelType string, userId, postId int) (bool, error) {
	key := modelType + fmt.Sprintf("%d", userId)
	flag, err := rd.rdb.GetBit(ctx, key, int64(postId)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	} else if flag == 0 {
		return false, nil
	}

	return true, nil
}

// IsUserLikePost is user like the post
func (rd *RedisDao) IsUserLikePost(userId, postId int) (bool, error) {
	return rd.isModelPost(model.IsLikePrefix, userId, postId)
}

// IsUserStarPost is user star the post
func (rd *RedisDao) IsUserStarPost(userId, postId int) (bool, error) {
	return rd.isModelPost(model.IsStarPrefix, userId, postId)
}

// IsUserFollow is user follow the other
func (rd *RedisDao) IsUserFollow(followerId, followeeId int) (bool, error) {
	return rd.isModelPost(model.IsFollowPrefix, followerId, followeeId)
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

func (rd *RedisDao) delCommon(postType, userType, isType string, userId, postId int) error {
	pipe := rd.rdb.Pipeline()

	// remove from the post's like list
	key := postType + fmt.Sprintf("%d", postId)
	pipe.SRem(ctx, key, userId).Result()
	// remove from the user's like list
	key = userType + fmt.Sprintf("%d", userId)
	pipe.LRem(ctx, key, 1, postId)
	// set like flag to false
	key = isType + fmt.Sprintf("%d", userId)
	pipe.SetBit(ctx, key, int64(postId), 0)

	_, err := pipe.Exec(ctx)

	return err
}

// DelLike dislike the post
func (rd *RedisDao) DelLike(userId, postId int) error {
	return rd.delCommon(model.PostLikePrefix, model.UserLikePrefix, model.IsLikePrefix, userId, postId)
}

// DelStar un star the post
func (rd *RedisDao) DelStar(userId, postId int) error {
	return rd.delCommon(model.PostStarPrefix, model.UserStarPrefix, model.IsStarPrefix, userId, postId)
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

// GetSessionUser session resolve
func (rd *RedisDao) GetSessionUser(sessionId string) (*model.User, error) {
	key := model.SessionPrefix + sessionId
	uid, err := rd.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("user not login")
	} else if err != nil {
		return nil, nil
	}

	userId, err := strconv.Atoi(uid)
	if err != nil {
		return nil, err
	}

	// check if session_id match
	user, err := rd.GetUserByID(userId)
	if user.SessionId != sessionId {
		return nil, errors.New("user not login")
	}

	return user, err
}

// SetSessionUser set session user
func (rd *RedisDao) SetSessionUser(sessionId string, userId, expire int) error {
	pipe := rd.rdb.Pipeline()

	key := model.SessionPrefix + sessionId
	pipe.Set(ctx, key, userId, time.Duration(expire)*time.Second)

	key = model.UinPrefix + fmt.Sprintf("%d", userId)
	pipe.HSet(ctx, key, "session_id", sessionId)

	_, err := pipe.Exec(ctx)
	return err
}

// DelSession delete the session
func (rd *RedisDao) DelSession(sessionId string) error {
	key := model.SessionPrefix + sessionId
	_, err := rd.rdb.Del(ctx, key).Result()

	return err
}

// GetHotPost hot post recommended to user
func (rd *RedisDao) getPostByIdList(idList []string) ([]*model.Post, error) {
	var postList []*model.Post
	for _, postId := range idList {
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

// GetHotPost hot post recommended to user
func (rd *RedisDao) GetHotPost(userId, count int) ([]*model.Post, error) {
	idList, err := rd.rdb.ZRandMember(ctx, model.AllPostModel, count, false).Result()
	if err != nil {
		return nil, err
	}

	return rd.getPostByIdList(idList)
}
