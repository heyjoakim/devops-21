package services

import (
	"fmt"

	"github.com/heyjoakim/devops-21/models"
)

var d = GetDbInstance()

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	err := d.db.First(&user, "username = ?", username).Error
	if err != nil {
		fmt.Println(err)
	}
	return user.UserID, err
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	err := d.db.First(&user, "user_id = ?", userID).Error
	if err != nil {
		fmt.Println(err)

	}
	return user
}

func CreateUser(user models.User) error {
	err := d.db.Create(&user).Error
	return err
}
