package models

type MessageDto struct {
	Content  string `gorm:"column:text" json:"content"`
	PubDate  int64  `gorm:"column:pub_date" json:"pub_date"`
	Username string `gorm:"column:username" json:"user"`
	Email    string `gorm:"column:email" json:"email"`
}
