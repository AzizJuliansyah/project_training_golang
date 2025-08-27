package entities

type Auth struct {
	Id			int
	Name		string
	Email		string
	Password	string
}

type Login struct{
	Email		string `validate:"required" label:"Email"`
	Password	string `validate:"required,min=5" label:"Password"`
}

type Register struct{
	Name		string `validate:"required" label:"Nama"`
	Email		string `validate:"required,email,isunique=users-email" label:"Email"`
	Password	string `validate:"required,min=5" label:"Password"`
	ConfirmPassword	string `validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}