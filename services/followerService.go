package services

import (
	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

func CreateFollower(follower models.Follower) error {
	err := GetDBInstance().db.Create(&follower).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err":      err,
			"follower": follower,
		}).Error("CreateFollower: DB err")
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
		log.WithFields(log.Fields{
			"err":              err,
			"followingUsersID": followingUsersID,
			"userToUnfollowID": userToUnfollowID,
		}).Error("UnfollowUser: DB err")
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
		log.WithFields(log.Fields{
			"err":         err,
			"userID":      userID,
			"noFollowers": noFollowers,
		}).Error("UnfollowUser: DB err")
	}
	return users
}

func GetUsersFollowedBy(userID uint) []models.Follower {
	var followers []models.Follower
	err := GetDBInstance().db.Where("who_id = ?", userID).Find(&followers).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err":    err,
			"userID": userID,
		}).Error("GetUsersFollowedBy: DB err")
	}
	return followers
}

func IsUserFollower(userID uint, followedID uint) bool {
	var follower models.Follower
	err := GetDBInstance().db.Where("who_id = ?", userID).
		Where("whom_id = ?", followedID).
		Find(&follower).Error

	if err != nil {
		log.WithFields(log.Fields{
			"err":        err,
			"userID":     userID,
			"followedID": followedID,
		}).Error("IsUserFollower: DB err")
	}
	return follower.WhoID != 0
}
