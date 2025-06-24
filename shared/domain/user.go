package domain

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type User struct {
	ID       string
	Name     string
	Barcode  string
	Role     Role
	Password string
}
