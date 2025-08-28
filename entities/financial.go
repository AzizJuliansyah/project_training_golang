package entities

import (
	"time"
)

type Financial struct {
	Id      	int64
	Date		time.Time
	Type		string
	Nominal		int64	
	Category	string
	Description	*string
	Attachment	*string
	UpdatedAt	time.Time
	CreatedAt	time.Time
}

type AddFinancial struct {
	Id			int64
	UserId		int
	Date		string	`validate:"required" label:"Tanggal"`
	Type		string		`validate:"required" label:"Tipe"`
	Nominal		int64		`validate:"required,numeric,min=0"`
	Category	string		`validate:"required" label:"Kategori"`
	Description	*string
	Attachment	*string
}