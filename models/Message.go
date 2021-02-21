package models

// Message model structure
type Message struct {
	MessageID uint   `gorm:"primaryKey;column:message_id"`
	AuthorID  uint   `gorm:"column:author_id"`
	Text      string `gorm:"column:text"`
	PubDate   int64  `gorm:"column:pub_date"`
	Flagged   int    `gorm:"column:flagged"`
}
