package api

import "github.com/gin-gonic/gin"

var GetUrls = map[string]gin.HandlerFunc{
	"/":                      Home,
	"/api/user/info":         GetUserInfo,
	"/api/user/logout":       Logout,
	"/api/user/myPost":       MyPost,
	"/api/user/myLike":       MyLike,
	"/api/user/myStar":       MyStar,
	"/api/user/myFollow":     MyFollow,
	"/api/user/hotPost":      HotPost,
	"/api/user/followerNum":  FollowerNum,
	"/api/user/followeeNum":  FolloweeNum,
	"/api/user/isLike":       IsLike,
	"/api/user/isStar":       IsStar,
	"/api/user/isFollow":     IsFollow,
	"/api/user/myFollower":   MyFollower,
	"/api/user/myFollowee":   MyFollowee,
	"/api/post/info":         GetPostInfo,
	"/api/post/likeNum":      GetLikeNum,
	"/api/post/starNum":      GetStarNum,
	"/api/post/commentNum":   GetCommentNum,
	"/api/post/comment/info": GetPostComment,
}

var PostUrls = map[string]gin.HandlerFunc{
	"/api/user/register":   Register,
	"/api/user/login":      Login,
	"/api/user/addFollow":  AddFollow,
	"/api/user/unFollow":   UnFollow,
	"/api/user/updateInfo": UpdateInfo,
	"/api/post/addPost":    AddPost,
	"/api/post/addComment": AddComment,
	"/api/post/addLike":    AddLike,
	"/api/post/addStar":    AddStar,
	"/api/post/delPost":    DelPost,
	"/api/post/delComment": DelComment,
	"/api/post/delLike":    DelLike,
	"/api/post/delStar":    DelStar,
}
