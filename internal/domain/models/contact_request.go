package models

type ContactRequest struct {
	Id         int `json:"id" bson:"_id"`
	FromUserId int `json:"fromUserId" bson:"fromuserid"`
	ToUserId   int `json:"toUserId" bson:"touserid"`
}
