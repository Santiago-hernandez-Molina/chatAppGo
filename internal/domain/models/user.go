package models

type User struct {
	Id       int    `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Status   bool   `json:"status" bson:"status"`
	Code     int    `json:"code" bson:"code"`
}

type UserWithToken struct {
	User  *User
	Token string
}
