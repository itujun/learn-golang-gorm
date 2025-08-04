package learn_golang_gorm

import "time"

type Address struct {
	ID        int64		`gorm:"primary_key;autoIncrement;column:id"`
	UserId    string	`gorm:"column:user_id"`
	Address   string	`gorm:"column:address"`
	CreatedAt time.Time	`gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time	`gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (a *Address) TableName() string {
	return "addresses"
}