package api

import (
	"github.com/gin-gonic/gin"
	"github.com/heroyan/twitter/model"
	"net/http"
	"strconv"
)

func GetPostInfo(c *gin.Context) {
	id := c.Query("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	svc := getPostSvc()

	post, err := svc.GetPost(postId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "post not found"})
		return
	}

	isLike := false
	isStar := false
	//if logged, then display if i like and star this post
	user, _ := checkLogin(c)
	if user != nil {
		isLike, _ = svc.IsUserLikePost(user.Id, postId)
		isStar, _ = svc.IsUserStarPost(user.Id, postId)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": post,
		"isLike": isLike, "isStar": isStar,
	})
}

func getLikeNum(c *gin.Context) {
	id := c.Query("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	num, _ := getPostSvc().GetLikeNum(postId)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
}

func getStarNum(c *gin.Context) {
	id := c.Query("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	num, _ := getPostSvc().GetStarNum(postId)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
}

func getCommentNum(c *gin.Context) {
	id := c.Query("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	num, _ := getPostSvc().GetCommentNum(postId)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": num})
}

func GetPostComment(c *gin.Context) {
	id := c.Query("id")
	postId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	start, size := getPagination(c)
	cmtList, err := getPostSvc().GetCommentByPost(postId, start, size)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cmtList})
}

func AddPost(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}

	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	if !checkPostContent(post.Content) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "content length invalid"})
		return
	}
	post.UserId = user.Id

	err := getPostSvc().AddPost(&post)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": post.Id})
}

func AddComment(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	var cmt model.Comment
	if err := c.ShouldBindJSON(&cmt); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	cmt.UserId = user.Id
	if !checkPostComment(&cmt) {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "content length or postId invalid"})
		return
	}

	err := getPostSvc().AddComment(&cmt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": cmt.Id})
}

func DelPost(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	// only can delete oneself post
	var post model.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	err := getPostSvc().DelPost(user.Id, post.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func DelComment(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	// only can delete oneself comment
	var cmt model.Comment
	if err := c.ShouldBindJSON(cmt); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	err := getPostSvc().DelComment(cmt.Id, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func AddLike(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	var like model.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	like.UserId = user.Id
	err := getPostSvc().AddLike(&like)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func AddStar(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	var star model.Star
	if err := c.ShouldBindJSON(&star); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	star.UserId = user.Id
	err := getPostSvc().AddStar(&star)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func DelLike(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	var like model.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	err := getPostSvc().DelLike(user.Id, like.PostId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}

func DelStar(c *gin.Context) {
	user, isLogin := checkLogin(c)
	if !isLogin {
		return
	}
	var star model.Star
	if err := c.ShouldBindJSON(&star); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}
	err := getPostSvc().DelStar(user.Id, star.PostId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ""})
}
