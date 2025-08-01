package learn_golang_gorm

import "time"

// Karena sudah sesuai dengan konvensi GORM,
// kita tidak perlu mendefinisikan nama tabel dan juga nama kolom (opsional).
// GORM akan secara otomatis menggunakan nama struct sebagai nama tabel dan nama field sebagai nama kolom
type User struct {
	ID       string 
	Name     Name `gorm:"embedded"`
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Name struct{
	FirstName 	string
	MiddleName 	string
	LastName  	string
}

// jika ingin mendefinisikan nama tabel
// func (u *User) TableName() string {
// 	return "users"
// }

type UserLog struct {
	ID 			int 	`gorm:"primary_key;autoIncrement"`
	UserId 		string	`gorm:"column:user_id"`
	Action 		string
	CreatedAt time.Time
	UpdatedAt time.Time
}