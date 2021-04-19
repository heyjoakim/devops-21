package services

import (
	"github.com/heyjoakim/devops-21/models"
)

func CreateFollower(follower models.Follower) error {
	err := GetDBInstance().db.Create(&follower).Error
	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data:    follower,
		})
	}
	return err
}

func UnfollowUser(followingUsersID uint, userToUnfollowID uint) error {
	var follower models.Follower
	err := GetDBInstance().db.Where("who_id = ?", followingUsersID).
		Where("whom_id = ?", userToUnfollowID).
		Delete(&follower).
		Error
	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data: map[string]uint{
				"followingUsersID": followingUsersID,
				"userToUnfollowID": userToUnfollowID,
			},
		})
	}
	return err
}

func GetAllUsersFollowers(userID uint, noFollowers int) []string {
	var users []string
	err := GetDBInstance().db.Model(&models.User{}).
		Select("\"user\".username").
		Joins("LEFT JOIN follower ON (follower.whom_id = \"user\".user_id)").
		Where("follower.who_id=?", userID).
		Limit(noFollowers).
		Scan(&users).Error

	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data: map[string]interface{}{
				"userID":      userID,
				"noFollowers": noFollowers,
			},
		})
	}
	return users
}

func GetUsersFollowedBy(userID uint) []models.Follower {
	var followers []models.Follower
	err := GetDBInstance().db.Where("who_id = ?", userID).Find(&followers).Error
	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data:    userID,
		})
	}
	return followers
}

func IsUserFollower(userID uint, followedID uint) bool {
	var follower models.Follower
	err := GetDBInstance().db.Where("who_id = ?", userID).
		Where("whom_id = ?", followedID).
		Find(&follower).Error

	if err != nil {
		LogError(models.Log{
			Message: err.Error(),
			Data: map[string]interface{}{
				"userID":     userID,
				"followedID": followedID,
			},
		})
	}
	return follower.WhoID != 0
}
