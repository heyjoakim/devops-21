package models

// MessageResponse defines a message reponse
type MessageResponse struct {
	Content string `gorm:"column:text" json:"content"`
	PubDate int    `gorm:"column:pub_date" json:"pub_date"`
	User    string `gorm:"column:username" json:"user"`
}
