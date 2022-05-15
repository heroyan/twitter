package api

import "github.com/gin-gonic/gin"

var GetUrls = map[string]gin.HandlerFunc{
	"/":                      Home,
	"/api/user/info":         GetUserInfo,
	"/api/user/logout":       Logout,
	"/api/user/myPost":       MyPost,
	"/api/user/myLike":       MyLike,
	"/api/user/myStar":       MyStar,
	"/api/post/info":         GetPostInfo,
	"/api/post/likeNum":      getLikeNum,
	"/api/post/starNum":      getStarNum,
	"/api/post/commentNum":   getCommentNum,
	"/api/post/comment/info": GetPostComment,
}

var PostUrls = map[string]gin.HandlerFunc{
	"/api/user/register":   Register,
	"/api/user/login":      Login,
	"/api/post/addPost":    AddPost,
	"/api/post/addComment": AddComment,
	"/api/post/addLike":    AddLike,
	"/api/post/addStar":    AddStar,
	"/api/post/delPost":    DelPost,
	"/api/post/delComment": DelComment,
	"/api/post/delLike":    DelLike,
	"/api/post/delStar":    DelStar,
}
