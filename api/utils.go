package api

import (
	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/config"
	"github.com/heroyan/twitter/dao"
	"github.com/heroyan/twitter/model"
	"github.com/heroyan/twitter/service"
	uuid "github.com/satori/go.uuid"
	"regexp"
)

func getUserSvc() *service.UserService {
	daoObj := dao.NewRedisDao(config.GetAddr(), config.GetPasswd(), config.GetDB())
	return service.NewUserService(daoObj)
}

// checkPassword 检查密码长度和组成
func checkPassword(password string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-z_$@A-Z0-9]{6,16}$", password); !ok {
		return false
	}
	return true
}

// checkUsername 检查user_name是否符合规则，长度不能超过64，只能有
func checkUsername(username string) (b bool) {
	if ok, _ := regexp.MatchString("^[a-z_A-Z0-9]{4,32}$", username); !ok {
		return false
	}
	return true
}

func genSessionId() string {
	u1 := uuid.Must(uuid.NewV4(), nil)
	return u1.String()
}

func getSessionUser(c *gin.Context) (*model.User, error) {
	sessionId, err := c.Cookie(config.GetSessionKey())

	if err != nil {
		return nil, err
	}

	user, err := getUserSvc().GetSessionUser(sessionId)

	return user, err
}
