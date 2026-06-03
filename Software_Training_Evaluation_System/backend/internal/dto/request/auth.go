package request

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	RealName string `json:"real_name" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=teacher student"`
}
