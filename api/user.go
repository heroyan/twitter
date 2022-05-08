package api

import (
	"github.com/heroyan/twitter/config"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
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
	c.Redirect(http.StatusTemporaryRedirect, "/")
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
