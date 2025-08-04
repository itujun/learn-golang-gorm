package learn_golang_gorm

import (
	"time"

	"gorm.io/gorm"
)

// model todos
type Todo struct {
	ID          int64			`gorm:"column:id;primary_key;autoIncrement"`
	UserId      string			`gorm:"column:user_id"`
	Title       string			`gorm:"column:title"`
	Description string			`gorm:"column:description"`
	Created_at  time.Time		`gorm:"column:created_at;autoCreateTime"`
	Updated_at  time.Time		`gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Deleted_at	gorm.DeletedAt 	`gorm:"column:deleted_at"`
}

// func (t *Todo) TableName() string {
// 	return "todos"
// }