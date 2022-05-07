package model

type User struct {
	Id            int    `json:"id" redis:"id"`
	UserName      string `json:"user_name" redis:"user_name"`
	Passwd        string `json:"passwd" redis:"passwd"`
	Nick          string `json:"nick" redis:"nick"`
	Name          string `json:"name" redis:"name"`
	Gender        bool   `json:"gender" redis:"gender"`
	Age           int    `json:"age" redis:"age"`
	SessionId     string `json:"session_id" redis:"session_id"`
	LastLoginTime int    `json:"last_login_time" redis:"last_login_time"`
	CreateTime    int    `json:"create_time" redis:"create_time"`
}

type Follow struct {
	FollowerId int `json:"follower_id" redis:"follower_id"`
	FolloweeId int `json:"followee_id" redis:"followee_id"`
	CreateTime int `json:"create_time" redis:"create_time"`
}
