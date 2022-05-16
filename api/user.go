package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/config"
	"github.com/heroyan/twitter/model"
)

func Home(c *gin.Context) {
	c.String(http.StatusOK, "home page")
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	if !checkUsername(user.UserName) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "user_name invalid"})
		return
	}

	if !checkPassword(user.Passwd) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "passwd invalid"})
		return
	}

	_, err := getUserSvc().RegisterUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	// if already login, then do nothing
	user2, _ := getSessionUser(c)
	if user2 != nil && user2.UserName == user.UserName {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "already login"})
		return
	}

	// 生成session_id到cookie
	err := getUserSvc().LoginUser(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	sessionId := genSessionId()
	c.SetCookie(config.GetSessionKey(), sessionId, config.GetSessionAge(),
		"/", "localhost", false, true)
	getUserSvc().SetSessionUser(sessionId, user.Id, config.GetSessionAge())

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func Logout(c *gin.Context) {
	sessionId, err := c.Cookie(config.GetSessionKey())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	err = getUserSvc().Logout(sessionId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.SetCookie(config.GetSessionKey(), sessionId, -1,
		"/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func GetUserInfo(c *gin.Context) {
	user, err := getSessionUser(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{
		"id":        user.Id,
		"user_name": user.UserName,
		"nick":      user.Nick,
		"name":      user.Name,
	}})
}

// MyPost posted by myself
func MyPost(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	start, size := getPagination(c)
	postList, err := getUserSvc().GetPostByUser(user.Id, start, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": postList})
}

// MyStar posts stared by myself
func MyStar(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	start, size := getPagination(c)
	postList, err := getUserSvc().GetPostStarByUser(user.Id, start, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": postList})
}

func MyFollow(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	start, size := getPagination(c)
	postList, err := getUserSvc().GetPostFollowByUser(user.Id, start, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": postList})
}

// MyLike posts liked by myself
func MyLike(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	start, size := getPagination(c)
	postList, err := getUserSvc().GetPostLikeByUser(user.Id, start, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": postList})
}

func HotPost(c *gin.Context) {
	user, _ := getSessionUser(c)
	userId := 0
	if user != nil {
		userId = user.Id
	}

	_, size := getPagination(c)
	postList, err := getUserSvc().GetHotPost(userId, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": postList})
}

func FollowerNum(c *gin.Context) {
	id := c.Query("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		// get current user's data, ignore errors
		user, _ := getSessionUser(c)
		if user != nil {
			userId = user.Id
		}
	}

	// ignore errors
	num, _ := getUserSvc().GetFollowerNum(userId)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
}

func FolloweeNum(c *gin.Context) {
	id := c.Query("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		// get current user's data, ignore errors
		user, _ := getSessionUser(c)
		if user != nil {
			userId = user.Id
		}
	}

	// ignore errors
	num, _ := getUserSvc().GetFolloweeNum(userId)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
}

func IsLike(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var postIdList []int
	idList := c.QueryArray("idList[]")
	for _, id := range idList {
		postId, _ := strconv.Atoi(id)
		postIdList = append(postIdList, postId)
	}

	likes, err := getUserSvc().IsLike(user.Id, postIdList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": likes})
}

func IsStar(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var postIdList []int
	idList := c.QueryArray("idList[]")
	for _, id := range idList {
		postId, _ := strconv.Atoi(id)
		postIdList = append(postIdList, postId)
	}

	stars, err := getUserSvc().IsStar(user.Id, postIdList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": stars})
}

func IsFollow(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var userIdList []int
	idList := c.QueryArray("idList[]")
	for _, id := range idList {
		uid, _ := strconv.Atoi(id)
		userIdList = append(userIdList, uid)
	}

	follows, err := getUserSvc().IsFollow(user.Id, userIdList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": follows})
}

func MyFollower(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	// get rand 100 ones
	users, err := getUserSvc().GetFollowers(user.Id, 100)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
}

func MyFollowee(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	// get rand 100 ones
	users, err := getUserSvc().GetFollowees(user.Id, 100)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": users})
}

func AddFollow(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var user2 model.User
	if err := c.ShouldBindJSON(&user2); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	err := getUserSvc().AddFollower(&model.Follow{
		FollowerId: user.Id,
		FolloweeId: user2.Id,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func UnFollow(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var user2 model.User
	if err := c.ShouldBindJSON(&user2); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	err := getUserSvc().UnFollow(&model.Follow{
		FollowerId: user.Id,
		FolloweeId: user2.Id,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}
