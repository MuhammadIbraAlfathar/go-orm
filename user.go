package go_orm

import "time"

type User struct {
	ID        string    `gorm:"primary_key;column:id"`
	Name      string    `gorm:"column:name"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateTime;autoUpdateTime"`
}
