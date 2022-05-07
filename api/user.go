package api

import (
	"github.com/heroyan/twitter/config"
	"github.com/heroyan/twitter/dao"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/model"
	"github.com/heroyan/twitter/service"
)

func getUserSvc() *service.UserService {
	daoObj := dao.NewRedisDao(config.GetAddr(), config.GetPasswd(), config.GetDB())
	return service.NewUserService(daoObj)
}

func Register(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok"})
}

func Logout(c *gin.Context) {

}
