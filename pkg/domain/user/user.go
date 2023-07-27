package domain

import "time"

type User struct {
	ID        string    `gorm:"not null;primary_key"`
	Email     string    `gorm:"unique;not null;column:email"`
	Password  string    `gorm:"not null;column:password"`
	CreatedAt time.Time `gorm:"not null;autoCreateTime;column:createdAt"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime;column:updatedAt"`
}

func (User) TableName() string {
	return "users"
}
