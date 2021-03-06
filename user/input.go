package user

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailAvailabilityInput struct {
	Email string `json:"email" binding:"required,email"`
}

type FormCreateInput struct {
	Name       string `form:"name" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Occupation string `form:"occupation" binding:"required"`
	Password   string `form:"password" binding:"required"`
	Error      error
}

type FormUpdateInput struct {
	ID         int
	Name       string `form:"name" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Occupation string `form:"occupation" binding:"required"`
	Error      error
}

type FormUpdateAvatar struct {
	ID             int
	Name           string `form:"name" binding:"required"`
	AvatarFileName string `form:"avatar"`
	Error          error
}

type WebLoginInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}
