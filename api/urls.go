package api

import "github.com/gin-gonic/gin"

var GetUrls = map[string]gin.HandlerFunc{
	"/user/info":         Register,
	"/post/info":         Login,
	"/post/comment/info": Logout,
	"/user/logout":       Logout,
}

var PostUrls = map[string]gin.HandlerFunc{
	"/user/register": Register,
	"/user/login":    Login,
}
