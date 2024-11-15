package entity

import "time"

type User struct {
	ID        int64     `pg:"id,pk"`
	Name      string    `pg:"name,notnull"`
	Email     string    `pg:"email,notnull"`
	CreatedAt time.Time `pg:"created_at"`
	UpdatedAt time.Time `pg:"updated_at"`
}

type CreateUserInput struct {
	Name  string `pg:"name"`
	Email string `pg:"email"`
}

type UpdateUserInput struct {
	ID    int64  `pg:"id"`
	Name  string `pg:"name"`
	Email string `pg:"email"`
}
