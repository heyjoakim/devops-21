package services

import (
	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	getUserIDErr := GetDBInstance().db.First(&user, "username = ?", username).Error
	if getUserIDErr != nil {
		log.WithFields(log.Fields{
			"err":      getUserIDErr,
			"username": username,
		}).Error("Error in GetUserID")
	}
	return user.UserID, getUserIDErr
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := GetDBInstance().db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"username": username,
		}).Error("GetUserFromUsername error")
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	getUserErr := GetDBInstance().db.First(&user, "user_id = ?", userID).Error
	if getUserErr != nil {
		log.WithFields(log.Fields{
			"getUserErr": getUserErr,
			"userID":     userID,
		}).Error("GetUser error")
	}
	return user
}

// CreateUser creates a new user in the database
func CreateUser(user models.User) error {
	createUserErr := GetDBInstance().db.Create(&user).Error
	if createUserErr != nil {
		log.WithFields(log.Fields{
			"createUserErr": createUserErr,
			"userObject":    user,
		}).Error("CreateUser error")
	}
	return createUserErr
}

// GetUserCount returns the number of users reigstered in the system
func GetUserCount() int64 {
	var count int64
	countErr := GetDBInstance().db.Find(&models.User{}).Count(&count).Error
	if countErr != nil {
		log.WithFields(log.Fields{
			"GetUserCountErr": countErr,
		}).Error("GetUserCount: DB err")
	}
	return count
}
