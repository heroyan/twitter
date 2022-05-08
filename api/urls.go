package api

import "github.com/gin-gonic/gin"

var GetUrls = map[string]gin.HandlerFunc{
	"/":                      Home,
	"/api/user/info":         GetUserInfo,
	"/api/user/logout":       Logout,
	"/api/post/info":         GetPostInfo,
	"/api/post/comment/info": GetPostComment,
}

var PostUrls = map[string]gin.HandlerFunc{
	"/api/user/register": Register,
	"/api/user/login":    Login,
}
