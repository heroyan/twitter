package model

type Post struct {
	Id         int    `json:"id" redis:"id"`
	Title      string `json:"title" redis:"title"`
	Content    string `json:"content" redis:"content"`
	UserId     int    `json:"user_id" redis:"user_id"`
	CreateTime int    `json:"create_time" redis:"create_time"`
}

type Like struct {
	Id         int `json:"id" redis:"id"`
	PostId     int `json:"post_id" redis:"post_id"`
	UserId     int `json:"user_id" redis:"user_id"`
	CreateTime int `json:"create_time" redis:"create_time"`
}

type Comment struct {
	Id         int    `json:"id" redis:"id"`
	PostId     int    `json:"post_id" redis:"post_id"`
	UserId     int    `json:"user_id" redis:"user_id"`
	Content    string `json:"content" redis:"content"`
	CreateTime int    `json:"create_time" redis:"create_time"`
}

type Star struct {
	Id         int `json:"id" redis:"id"`
	PostId     int `json:"post_id" redis:"post_id"`
	UserId     int `json:"user_id" redis:"user_id"`
	CreateTime int `json:"create_time" redis:"create_time"`
}
