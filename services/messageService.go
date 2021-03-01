package services

import "github.com/heyjoakim/devops-21/models"

func GetPublicMessages(numberOfMessages int) []models.MessageDto {
	var results []models.MessageDto
	d.db.Model(&models.Message{}).
		Select("message.text, message.pub_date, user.email, user.username").
		Joins("left join user on user.user_id = message.author_id").
		Where("message.flagged=0").
		Order("pub_date desc").
		Limit(numberOfMessages).
		Scan(&results)

	return results
}

func GetMessagesForUser(numberOfMessages int, userId uint) []models.MessageDto {
	var results []models.MessageDto
	d.db.Model(models.Message{}).
		Order("pub_date desc").
		Select("message.text,message.pub_date, user.email, user.username").
		Joins("left join user on user.user_id = message.author_id").
		Where("user.user_id=?", userId).
		Limit(numberOfMessages).
		Scan(&results)
	return results
}

func CreateMessage(message models.Message) error {
	err := d.db.Create(&message).Error
	return err
}
