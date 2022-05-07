package service

import (
	"fmt"
	"github.com/heroyan/twitter/model"
	"math/rand"
	"testing"
	"time"

	"github.com/heroyan/twitter/dao"
	"github.com/smartystreets/goconvey/convey"
)

var userSvc = NewUserService(dao.NewRedisDao("localhost:6379", "", 0))

func TestUserService_RegisterUser(t *testing.T) {
	convey.Convey("RegisterUser test", t, func() {
		rand.Seed(time.Now().Unix())
		user := &model.User{
			Id:            0,
			UserName:      fmt.Sprintf("fagenb_%d", rand.Intn(10000)),
			Passwd:        "fagezhenniu",
			Nick:          "fage",
			Name:          "shuifa",
			Gender:        false,
			Age:           18,
			LastLoginTime: 0,
			CreateTime:    0,
		}
		id, err := userSvc.RegisterUser(user)
		convey.So(err, convey.ShouldBeNil)
		convey.So(id, convey.ShouldNotBeZeroValue)

		// pass reset
		user.Passwd = "fagezhenniu"
		err = userSvc.LoginUser(user)
		convey.So(err, convey.ShouldBeNil)
	})
}
