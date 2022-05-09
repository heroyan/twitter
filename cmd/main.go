//+build linux darwin windows mac

package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/api"
	"github.com/heroyan/twitter/config"
	"time"
)

var configFile = flag.String("conf", "./config.json", "--conf=filepath")
var configType = flag.String("conf-type", "json", "conf-type=json")

func main() {
	err := config.LoadConfig(*configFile, *configType)
	if err != nil {
		panic("load config file:" + *configFile + ", type:" + *configType + ", error: " + err.Error())
	}
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// log format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	for url, handleFunc := range api.GetUrls {
		router.GET(url, handleFunc)
	}
	for url, handleFunc := range api.PostUrls {
		router.POST(url, handleFunc)
	}

	router.Run(":9981")
}
