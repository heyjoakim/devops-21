package models

// User model
type User struct {
	UserID   uint      `gorm:"primaryKey"`
	Username string    `gorm:"column:username"`
	Email    string    `gorm:"column:email"`
	PwHash   string    `gorm:"column:pw_hash"`
	Messages []Message `gorm:"foreignKey:AuthorID"`
}
