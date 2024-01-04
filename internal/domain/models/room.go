package models

type Room struct {
	Id    int        `json:"id" bson:"_id"`
	Name  string     `json:"name" bson:"name"`
	Users []UserRoom `json:"users" bson:"users"`
}
