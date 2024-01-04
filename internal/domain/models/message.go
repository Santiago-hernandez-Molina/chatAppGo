package models

type Message struct {
	Id      int    `json:"id" bson:"_id"`
	Content string `json:"content" bson:"content"`
	RoomId  int    `json:"roomId" bson:"roomid"`
	UserId  int    `json:"userId" bson:"userid"`
}

type MessageUser struct {
	Id      int    `json:"id" bson:"_id"`
	Content string `json:"content" bson:"content"`
	User    *User   `json:"user" bson:"user"`
}
