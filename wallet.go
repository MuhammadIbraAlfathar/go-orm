package go_orm

import "time"

type Wallet struct {
	ID        string    `gorm:"primary_key;column:id"`
	UserId    string    `gorm:"column:user_id"`
	Balance   int64     `gorm:"column:balance"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:created_at;autoCreateTime;autoUpdateTime"`
}

func (w *Wallet) TableName() string {
	return "wallets"
}
