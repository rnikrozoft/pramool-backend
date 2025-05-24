package entity

import "time"

type User struct {
	UserId    string    `db:"user_id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
