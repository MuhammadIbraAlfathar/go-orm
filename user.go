package go_orm

import "time"

type User struct {
	ID        string    `gorm:"primary_key;column:id;<-:create"`
	Name      Name      `gorm:"embedded"`
	Password  string    `gorm:"column:password"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateTime;autoUpdateTime"`
	Wallet    Wallet    `gorm:"foreignKey:user_id;references:id"`
	Address   []Address `gorm:"foreignKey:user_id;references:id"`
}

type Name struct {
	FirstName  string `gorm:"column:first_name"`
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name"`
}

type UserLog struct {
	ID        string    `gorm:"primary_key;column:id;autoIncrement"`
	UserID    string    `gorm:"column:user_id"`
	Action    string    `gorm:"column:action"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateTime;autoUpdateTime"`
}
