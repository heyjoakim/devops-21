package services

import "github.com/heyjoakim/devops-21/models"

// Returns x number of public messages that have not been flagged, in desc order by publish date
// "user" is a reserved word in postgres, so it needs to be quoted in the queries
func GetPublicMessages(numberOfMessages int) []models.MessageDto {
	var results []models.MessageDto
	d.DB.Model(&models.Message{}).
		Select("message.text, message.pub_date, \"user\".username, \"user\".email").
		Joins("left join \"user\" on message.author_id = \"user\".user_id").
		Where("message.flagged=0").
		Order("pub_date desc").
		Limit(numberOfMessages).
		Scan(&results)
	return results
}

// Returns x number of messages for the specified user that have not been flagged,
// in desc order by publish date
// "user" is a reserved word in postgres, so it needs to be quoted in the queries
func GetMessagesForUser(numberOfMessages int, userId uint) []models.MessageDto {
	var results []models.MessageDto
	d.DB.Model(models.Message{}).
		Order("pub_date desc").
		Select("message.text,message.pub_date, \"user\".email, \"user\".username").
		Joins("left join \"user\" on \"user\".user_id = message.author_id").
		Where("\"user\".user_id=?", userId).
		Limit(numberOfMessages).
		Scan(&results)

	return results
}

func CreateMessage(message models.Message) error {
	err := d.DB.Create(&message).Error
	return err
}
