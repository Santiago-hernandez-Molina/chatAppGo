package models

type UserRoom struct {
	UserId int    `json:"userId" bson:"userid"`
	RoomId int    `json:"roomid" bson:"roomid"`
	Role   string `json:"role" bson:"role"`
}
