package entity

type User struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserInput struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
