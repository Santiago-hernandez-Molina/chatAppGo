package models

type Room struct {
	Id    int        `json:"id" bson:"_id"`
	Name  string     `json:"name" bson:"name"`
	Users []UserRoom `json:"users" bson:"users"`
	Type  RoomType   `json:"roomType" bson:"roomtype"`
}

type RoomType int64

const (
	Contact int = 0
	Group   int = 1
)
