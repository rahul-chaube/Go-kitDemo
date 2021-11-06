package post

import "time"

type Post struct {
	Id       string    `json:"id"`
	User     string    `json:"user_id"`
	Message  string    `json:"message"`
	MediaUrl string    `json:"image"`
	Type     string    `json:"type"`
	Comment  []Comment `json:"address"`
}
type Comment struct {
	Id      string    `json:"id"`
	Text    string    `json:"message"`
	UserId  string    `json:"user_id"`
	Post_id string    `json:"post_id"`
	Time    time.Time `json:"time"`
}
