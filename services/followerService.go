package services

import "github.com/heyjoakim/devops-21/models"

func CreateFollower(follower models.Follower) error {
	return d.db.Create(&follower).Error
}

func UnfollowUser(followingUsersID uint, userToUnfollowID uint) error {
	var follower models.Follower
	err := d.db.Where("who_id = ?", followingUsersID).
		Where("whom_id = ?", userToUnfollowID).
		Delete(&follower).
		Error
	return err
}

func GetAllUsersFollowers(userID uint, noFollowers int) []string {
	var users []string
	d.db.Model(&models.User{}).
		Select("\"user\".username").
		Joins("LEFT JOIN follower ON (follower.whom_id = \"user\".user_id)").
		Where("follower.who_id=?", userID).
		Limit(noFollowers).
		Scan(&users)
	return users
}

func GetUsersFollowedBy(userID uint) []models.Follower {
	var followers []models.Follower
	d.db.Where("who_id = ?", userID).Find(&followers)
	return followers
}

func IsUserFollower(userID uint, followedID uint) bool {
	var follower models.Follower
	d.db.Where("who_id = ?", userID).
		Where("whom_id = ?", followedID).
		Find(&follower)
	return follower.WhoID != 0
}
