package dtos

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
    Id int `json:"id"`
    Email string `json:"email"`
    Username string `json:"username"`
}

type ActivationRequest struct{
    Email string `json:"email"`
    Code int `json:"code"`
}
