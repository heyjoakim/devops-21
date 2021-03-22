package services

import (
	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

var d = GetDBInstance()

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	err := d.db.First(&user, "username = ?", username).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"username": username,
		}).Error("Error in GetUserID")
	}
	return user.UserID, err
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := d.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"username": username,
		}).Error("GetUserFromUsername error")
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	err := d.db.First(&user, "user_id = ?", userID).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"userID": userID,
		}).Error("GetUser error")
	}
	return user
}

// CreateUser creates a new user in the database
func CreateUser(user models.User) error {
	err := d.db.Create(&user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"userObject": user,
		}).Error("CreateUser error")
	}
	return err
}

// GetUserCount returns the number of users reigstered in the system
func GetUserCount() int64 {
	var count int64
	err := d.db.Find(&models.User{}).Count(&count).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("GetUserCount: DB err")
	}
	return count
}
