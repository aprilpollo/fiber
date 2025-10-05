package models

type User struct {
	BaseModel

	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRole struct {
	BaseModel

	Name string `json:"name"`
	
}

func (UserRole) TableName() string {
	return "user_roles"
}

func (User) TableName() string {
	return "users"
}
