package model

type Post struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	UserId     int    `json:"user_id"`
	CreateTime int    `json:"create_time"`
}

type Like struct {
	Id         int `json:"id"`
	PostId     int `json:"post_id"`
	UserId     int `json:"user_id"`
	CreateTime int `json:"create_time"`
}

type Comment struct {
	Id         int    `json:"id"`
	PostId     int    `json:"post_id"`
	UserId     int    `json:"user_id"`
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
}

type Star struct {
	Id         int `json:"id"`
	PostId     int `json:"post_id"`
	UserId     int `json:"user_id"`
	CreateTime int `json:"create_time"`
}
