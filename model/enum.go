package model

const (
	UserModel    string = "user"
	PostModel           = "post"
	LikeModel           = "like"
	CommentModel        = "comment"
	AllPostModel        = "allPost"
)

const (
	UserPrefix          string = "user:"
	UserPostPrefix             = "user:post:"
	UserFollowerPrefix         = "user:follower:"
	UserFolloweePrefix         = "user:followee:"
	UserStarPrefix             = "user:star:"
	IsStarPrefix               = "user:is_star:"
	UserLikePrefix             = "user:like:"
	IsLikePrefix               = "user:is_like:"
	UserTimelinePrefix         = "user:timeline:"
	IsFollowPrefix             = "user:is_follow:"
	UinPrefix                  = "uin:"
	GenPrefix                  = "gen:"
	PostPrefix                 = "post:"
	PostLikePrefix             = "post:like:"
	PostCommentPrefix          = "post:comment:"
	PostStarPrefix             = "post:star:"
	CommentDetailPrefix        = "comment:detail:"
	SessionPrefix              = "session:"
)

const (
	NeedLoginCode int = 50008
)
