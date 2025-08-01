// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"database/sql"
	"time"
)

type Profile struct {
	ID        int32          `db:"id" json:"id"`
	UserID    int32          `db:"user_id" json:"user_id"`
	FirstName string         `db:"first_name" json:"first_name"`
	LastName  string         `db:"last_name" json:"last_name"`
	Age       sql.NullInt32  `db:"age" json:"age"`
	Gender    sql.NullString `db:"gender" json:"gender"`
	City      sql.NullString `db:"city" json:"city"`
	Interests []string       `db:"interests" json:"interests"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID           int32     `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"password_hash"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
