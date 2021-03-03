package services

import "github.com/heyjoakim/devops-21/models"

func CreateFollower(follower models.Follower) error {
	return d.DB.Create(&follower).Error
}

func UnfollowUser(followingUsersId uint, userToUnfollowId uint) error {
	var follower models.Follower
	err := d.DB.Where("who_id = ?", followingUsersId).
		Where("whom_id = ?", userToUnfollowId).
		Delete(&follower).
		Error
	return err
}

func GetAllUsersFollowers(userId uint, noFollowers int) []string {
	var users []string
	d.DB.Model(&models.User{}).
		Select("\"user\".username").
		Joins("LEFT JOIN follower ON (follower.whom_id = \"user\".user_id)").
		Where("follower.who_id=?", userId).
		Limit(noFollowers).
		Scan(&users)
	return users
}

func GetUsersFollowedBy(userId uint) []models.Follower {
	var followers []models.Follower
	d.DB.Where("who_id = ?", userId).Find(&followers)
	return followers
}

func IsUserFollower(userId uint, followedId uint) bool {
	var follower models.Follower
	d.DB.Where("who_id = ?", userId).
		Where("whom_id = ?", followedId).
		Find(&follower)
	if follower.WhoID != 0 {
		return true
	}
	return false
}
