package request

type UserQuery struct {
	Role    string `form:"role"`
	Keyword string `form:"keyword"`
}

type CreateUserReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin teacher student"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UpdateUserReq struct {
	RealName string `json:"real_name"`
	Password string `json:"password"`
	Role     string `json:"role" binding:"omitempty,oneof=admin teacher student"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   *int   `json:"status"`
}
