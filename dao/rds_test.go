package dao

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/heroyan/twitter/model"
	"github.com/smartystreets/goconvey/convey"
)

var rd = NewRedisDao("localhost:6379", "", 0)
var testUserId = 8
var testUserName = "yanshuifa8"

func TestNewRedisDao(t *testing.T) {
	user, err := rd.GetUser(testUserName)
	if err != nil {
		t.Errorf("get user error: %+v", err)
		return
	}
	if user == nil {
		t.Logf("%s not found", testUserName)
		return
	}
	t.Logf("found user: %+v", user)
}

func TestRedisDao_GenerateId(t *testing.T) {
	modelTypeList := []string{model.UserModel, model.LikeModel, model.PostModel, model.CommentModel}
	for _, md := range modelTypeList {
		uid, err := rd.generateId(md)
		if err != nil {
			t.Errorf("gen %s id err:%+v", md, err)
		}
		t.Logf("gen %s id: %+v", md, uid)
	}
}

func TestRedisDao_AddUser(t *testing.T) {
	convey.Convey("AddUser test", t, func() {
		user := &model.User{
			Id:            0,
			CreateTime:    int(time.Now().Unix()),
			Name:          "水发",
			Nick:          "发哥",
			UserName:      testUserName,
			Gender:        true,
			Age:           18,
			LastLoginTime: int(time.Now().Unix()),
			Passwd:        "pass",
		}
		user1, err := rd.GetUser(testUserName)
		t.Logf("get user2: %+v", user1)
		convey.So(err, convey.ShouldBeNil)
		if user1 != nil {
			t.Logf("user: %s exists", user1.UserName)
			return
		}

		err = rd.AddUser(user)
		convey.So(err, convey.ShouldBeNil)
		convey.So(user.Id, convey.ShouldNotBeZeroValue)

		user2, err := rd.GetUserByID(user.Id)
		t.Logf("get user2: %+v", user2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(user2, convey.ShouldNotBeNil)
		convey.So(user2.UserName, convey.ShouldEqual, testUserName)
	})
}

func TestRedisDao_AddPost(t *testing.T) {
	convey.Convey("AddPost test", t, func() {
		post := &model.Post{
			Id:         0,
			Title:      "test",
			Content:    "content is wonderful",
			UserId:     testUserId,
			CreateTime: int(time.Now().Unix()),
		}
		err := rd.AddPost(post)
		convey.So(err, convey.ShouldBeNil)
		t.Logf("new post id is: %d", post.Id)

		post2, err := rd.GetPost(post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(post2, convey.ShouldNotBeNil)
		t.Logf("get post info: %+v", post2)
	})
}

func addTestPost() (*model.Post, error) {
	rand.Seed(time.Now().Unix())
	post := &model.Post{
		Id:         0,
		Title:      fmt.Sprintf("title: %d", rand.Intn(10000)),
		Content:    "content is wonderful",
		UserId:     testUserId,
		CreateTime: int(time.Now().Unix()),
	}
	err := rd.AddPost(post)

	return post, err
}

func TestRedisDao_AddLike(t *testing.T) {
	convey.Convey("AddLike test", t, func() {
		post, err := addTestPost()
		convey.So(err, convey.ShouldBeNil)
		t.Logf("new post id is: %d", post.Id)
		convey.So(post.Id, convey.ShouldNotBeZeroValue)

		like := &model.Like{
			Id:         0,
			PostId:     post.Id,
			UserId:     testUserId,
			CreateTime: int(time.Now().Unix()),
		}
		err = rd.AddLike(like)
		convey.So(err, convey.ShouldBeNil)

		flag, err := rd.IsUserLikePost(testUserId, post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(flag, convey.ShouldBeTrue)
	})
}

func TestRedisDao_AddComment(t *testing.T) {
	convey.Convey("AddComment test", t, func() {
		post, err := addTestPost()
		convey.So(err, convey.ShouldBeNil)
		t.Logf("new post id is: %d", post.Id)
		convey.So(post.Id, convey.ShouldNotBeZeroValue)

		rand.Seed(time.Now().Unix())
		cmt := &model.Comment{
			Id:         0,
			PostId:     post.Id,
			UserId:     testUserId,
			Content:    fmt.Sprintf("content: %d", rand.Intn(10000)),
			CreateTime: int(time.Now().Unix()),
		}
		err = rd.AddComment(cmt)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cmt.Id, convey.ShouldNotBeZeroValue)

		cmt2, err := rd.GetCommentByID(cmt.Id)
		t.Logf("GetCommentByID: %d, %+v", cmt.Id, cmt2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(cmt2.Id, convey.ShouldEqual, cmt.Id)
	})
}

func TestRedisDao_AddStar(t *testing.T) {
	convey.Convey("AddStar test", t, func() {
		post, err := addTestPost()
		convey.So(err, convey.ShouldBeNil)
		t.Logf("new post id is: %d", post.Id)
		convey.So(post.Id, convey.ShouldNotBeZeroValue)

		star := &model.Star{
			Id:         0,
			PostId:     post.Id,
			UserId:     testUserId,
			CreateTime: int(time.Now().Unix()),
		}
		err = rd.AddStar(star)
		convey.So(err, convey.ShouldBeNil)

		flag, err := rd.IsUserStarPost(testUserId, post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(flag, convey.ShouldBeTrue)
	})
}

func TestRedisDao_AddFollower(t *testing.T) {
	convey.Convey("AddFollower test", t, func() {
		follow := &model.Follow{
			FollowerId: 10,
			FolloweeId: testUserId,
			CreateTime: int(time.Now().Unix()),
		}
		err := rd.AddFollower(follow)
		convey.So(err, convey.ShouldBeNil)

		num, err := rd.GetFollowerNum(testUserId)
		t.Logf("user %d has %d followers", testUserId, num)
		convey.So(err, convey.ShouldBeNil)
		convey.So(num, convey.ShouldNotBeZeroValue)

		follow.FollowerId = testUserId
		follow.FolloweeId = 10
		err = rd.AddFollower(follow)
		convey.So(err, convey.ShouldBeNil)

		num, err = rd.GetFolloweeNum(testUserId)
		t.Logf("user %d has followed %d users", testUserId, num)
		convey.So(err, convey.ShouldBeNil)
		convey.So(num, convey.ShouldNotBeZeroValue)
	})
}

func TestRedisDao_DelLike(t *testing.T) {
	convey.Convey("DelLike test", t, func() {
		post, err := addTestPost()
		convey.So(err, convey.ShouldBeNil)

		err = rd.AddLike(&model.Like{
			Id:         0,
			PostId:     post.Id,
			UserId:     testUserId,
			CreateTime: int(time.Now().Unix()),
		})
		convey.So(err, convey.ShouldBeNil)

		flag, err := rd.IsUserLikePost(testUserId, post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(flag, convey.ShouldBeTrue)

		err = rd.DelLike(testUserId, post.Id)
		convey.So(err, convey.ShouldBeNil)

		flag, err = rd.IsUserLikePost(testUserId, post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(flag, convey.ShouldBeFalse)
	})
}

func TestRedisDao_GetCommentByPost(t *testing.T) {
	convey.Convey("GetCommentByPost test", t, func() {
		post, err := addTestPost()
		convey.So(err, convey.ShouldBeNil)

		cmt := &model.Comment{
			Id:         0,
			PostId:     post.Id,
			UserId:     testUserId,
			Content:    "test content",
			CreateTime: 0,
		}
		err = rd.AddComment(cmt)
		convey.So(err, convey.ShouldBeNil)

		cmt.Content = "test content2"
		err = rd.AddComment(cmt)
		convey.So(err, convey.ShouldBeNil)

		cmtList, err := rd.GetCommentByPost(post.Id, 0, 100)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(cmtList), convey.ShouldEqual, 2)

		// to del a comment
		err = rd.DelComment(cmt.Id, post.Id)
		convey.So(err, convey.ShouldBeNil)

		num, err := rd.GetCommentNum(post.Id)
		convey.So(err, convey.ShouldBeNil)
		convey.So(num, convey.ShouldEqual, 1)
	})
}

func TestRedisDao_SetSessionUser(t *testing.T) {
	convey.Convey("SetSessionUser test", t, func() {
		session := "test-session-id"
		err := rd.SetSessionUser(session, 8, 86400)
		convey.So(err, convey.ShouldBeNil)
	})

}
