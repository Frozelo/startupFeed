package dto

type CreateUserDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateFeedbackDTO struct {
	UserId int64  `json:"user_id"`
	Text   string `json:"text"`
}

type LoginUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProjectDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
