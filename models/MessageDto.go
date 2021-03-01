package models

type MessageDto struct {
	Content  string `gorm:"column:text"`
	PubDate  int64  `gorm:"column:pub_date"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
}
