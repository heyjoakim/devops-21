package models

// MessageResponse defines a message reponse
type MessageResponse struct {
	Content string `json:"content"`
	PubDate int64  `json:"pub_date"`
	User    string `json:"user"`
}
