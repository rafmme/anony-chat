package domain

type User struct {
	ID       string `gorm:"primary_key"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}
