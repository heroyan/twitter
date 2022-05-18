package api

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/config"
	"github.com/heroyan/twitter/dao"
	"github.com/heroyan/twitter/model"
	"github.com/heroyan/twitter/service"
)

func getUserSvc() *service.UserService {
	daoObj := dao.NewRedisDao(config.GetAddr(), config.GetPasswd(), config.GetDB())
	return service.NewUserService(daoObj)
}

func getPostSvc() *service.PostService {
	daoObj := dao.NewRedisDao(config.GetAddr(), config.GetPasswd(), config.GetDB())
	return service.NewPostService(daoObj)
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
		// ignore error, mostly the session key not exist
		return nil, nil
	}

	user, err := getUserSvc().GetSessionUser(sessionId)

	return user, err
}

func checkPostContent(content string) bool {
	if content == "" || utf8.RuneCountInString(content) > 144 {
		return false
	}

	return true
}

func checkLogin(c *gin.Context, needLogin bool) (*model.User, bool) {
	user, err := getSessionUser(c)
	code := 0
	if needLogin {
		code = model.NeedLoginCode
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": code, "data": gin.H{}, "msg": err.Error()})
		return nil, false
	}
	if user == nil {
		c.JSON(http.StatusOK, gin.H{"code": code, "data": gin.H{}, "msg": "not login"})
		return nil, false
	}

	return user, true
}

func checkPostComment(comment *model.Comment) bool {
	if comment.PostId == 0 || comment.Content == "" || utf8.RuneCountInString(comment.Content) > 100 {
		return false
	}

	return true
}

func getPagination(c *gin.Context) (start, size int) {
	page := c.Query("page")
	limit := c.Query("limit")
	// start非法默认从0开始，忽略err
	pgNo, _ := strconv.Atoi(page)
	if pgNo < 1 {
		pgNo = 1
	}
	size, err := strconv.Atoi(limit)
	// 最大获取不能超过100
	if err != nil || size < 1 || size > 100 {
		size = 10
	}
	start = (pgNo - 1) * size

	return start, size
}
