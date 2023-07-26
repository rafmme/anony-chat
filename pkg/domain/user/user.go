package domain

type User struct {
	ID        string `gorm:"primary_key"`
	Email     string `gorm:"column:email"`
	Password  string `gorm:"column:password"`
	CreatedAt string `gorm:"column:createdAt"`
	UpdatedAt string `gorm:"column:updatedAt"`
}

func (User) TableName() string {
	return "users"
}
