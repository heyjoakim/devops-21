package models

type MessageViewModel struct {
	Content string `json:"content"`
	PubDate string `json:"pub_date"`
	User    string `json:"user"`
	Email   string `json:"email"`
}
