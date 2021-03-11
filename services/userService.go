package services

import (
	"log"

	"github.com/heyjoakim/devops-21/models"
)

var d = GetDbInstance()

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	err := d.db.First(&user, "username = ?", username).Error
	if err != nil {
		log.Println(err)
	}
	return user.UserID, err
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	err := d.db.First(&user, "user_id = ?", userID).Error
	if err != nil {
		log.Println(err)
	}
	return user
}

// CreateUser creates a new user in the database
func CreateUser(user models.User) error {
	err := d.db.Create(&user).Error
	return err
}

// GetUserCount returns the number of users reigstered in the system
func GetUserCount() int64 {
	var count int64
	d.db.Find(&models.User{}).Count(&count)
	return count
}
