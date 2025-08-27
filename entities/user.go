package entities

type User struct {
	Id			int
	Name		string `validate:"required" label:"Nama"`
	Email		string `validate:"required,email" label:"Email"`
	Password	string `validate:"omitempty,min=5"`
	Photo		*string
}