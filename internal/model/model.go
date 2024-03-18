package model

type UserCreateRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}

type UserListRequest struct {
	Limit int `query:"limit"`
	Page  int `query:"limit"`
}

type UserResponse struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
