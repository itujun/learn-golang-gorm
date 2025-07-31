package learn_golang_gorm

// Karena sudah sesuai dengan konvensi GORM,
// kita tidak perlu mendefinisikan nama tabel dan juga nama kolom (opsional).
// GORM akan secara otomatis menggunakan nama struct sebagai nama tabel dan nama field sebagai nama kolom
type User struct {
	ID       string
	Name     string
	Password string
	CreateAt string
	UpdateAt string
}

// jika ingin mendefinisikan nama tabel
// func (u *User) TableName() string {
// 	return "users"
// }