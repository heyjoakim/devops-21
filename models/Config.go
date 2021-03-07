package models

// Config model
type Config struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"column:key"`
	Value string `gorm:"column:value"`
}
