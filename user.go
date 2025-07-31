package learn_golang_gorm

// Karena sudah sesuai dengan konvensi GORM,
// kita tidak perlu mendefinisikan nama tabel dan juga nama kolom (opsional).
// GORM akan secara otomatis menggunakan nama struct sebagai nama tabel dan nama field sebagai nama kolom
type User struct {
	ID       string `gorm:"<-:create"` // <-:create berarti field ini hanya bisa diisi saat create, tidak bisa diupdate 
	Name     string
	Password string
	CreateAt string `gorm:"<-:create"` // <-:create berarti field ini hanya bisa diisi saat create, tidak bisa diupdate
	UpdateAt string
	Information string `gorm:"-"` // "-" berarti field ini tidak akan disimpan di database
}

// jika ingin mendefinisikan nama tabel
// func (u *User) TableName() string {
// 	return "users"
// }