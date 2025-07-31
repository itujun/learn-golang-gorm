package learn_golang_gorm

type User struct {
	ID       string `gorm:"primaryKey;column:id"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	CreateAt string `gorm:"column:created_at;autoCreateTime"`
	UpdateAt string `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

// jika ingin mendefinisikan nama tabel
func (u *User) TableName() string {
	return "users"
}