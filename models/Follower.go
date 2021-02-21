package models

// Follower model
type Follower struct {
	WhoID  uint `gorm:"column:who_id"`
	WhomID uint `gorm:"column:whom_id"`
}
