package models

type Message struct {
	Email    string `gorm:"column:user_id"`
	Username string `gorm:"column:user_id"`
	Text     string `gorm:"column:user_id"`
	PubDate  string `gorm:"column:user_id"`
}
