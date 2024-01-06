package models

type UserRoom struct {
	UserId int    `json:"userId" bson:"userid"`
	RoomId int    `json:"roomId" bson:"roomid"`
	Role   string `json:"role" bson:"role"`
}
