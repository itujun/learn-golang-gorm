package learn_golang_gorm

import "time"

// Karena sudah sesuai dengan konvensi GORM,
// kita tidak perlu mendefinisikan nama tabel dan juga nama kolom (opsional).
// GORM akan secara otomatis menggunakan nama struct sebagai nama tabel dan nama field sebagai nama kolom
type User struct {
	ID       string 
	Name     Name 		`gorm:"embedded"`
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
	Wallet   	Wallet 		`gorm:"foreignKey:user_id;references:id"`
	Addresses 	[]Address 	`gorm:"foreignKey:user_id;references:id"`
	LikeProducts []Product	`gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:product_id"`
	// ====================> `gorm:"many2many: nama_table_penghubung; foreignKey: nama_kolom_penghubung; joinForeignKey: nama_kolom_dari_tabel_1; references: nama_kolom_dari_tabel_1; joinReferences: nama_kolom_dari_tabel_2"` 
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
	CreatedAt 	int64	`gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt 	int64	`gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}