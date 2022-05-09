package config

import (
	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"

	"testing"
)

func TestGetSessionAge(t *testing.T) {
	convey.Convey("GetSessionAge test", t, func() {
		viper.SetConfigFile("/Users/shuifa/workspace/twitter/config.json")
		viper.SetConfigType("json")
		viper.ReadInConfig()

		t.Logf("%s", GetSessionKey())
		t.Logf("%d", GetSessionAge())
		t.Logf("%s", GetAddr())
	})
}
