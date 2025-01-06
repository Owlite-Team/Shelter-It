package request

type LoginReq struct {
	Email    string `json:"username" binding:"required, email"`
	Password string `json:"password" binding:"required, min=8"`
}
