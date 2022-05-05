package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/heroyan/twitter/model"
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

func (rd *RedisDao) GetUser(userName string) (user *model.User, err error) {
	key := "user:" + userName
	result, err := rd.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// user:username saves id in the redis
	// uin:id hash saves the user info
	user = &model.User{}
	key = "uin:" + result
	ret := rd.rdb.HGetAll(ctx, key)
	err = ret.Scan(user)

	return
}

func (rd *RedisDao) SaveUser(user *model.User) (err error) {
	return
}

func (rd *RedisDao) GetPost(id int) (post *model.Post, err error) {
	return
}

func (rd *RedisDao) SavePost(post *model.Post) (err error) {
	return
}
