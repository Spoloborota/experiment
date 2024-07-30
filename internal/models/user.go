package models

type User struct {
	ID           int
	FirstName    string
	LastName     string
	Age          int
	Gender       string
	Interests    string
	City         string
	PasswordHash string
}
