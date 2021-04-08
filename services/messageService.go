package services

import (
	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

// GetPublicMessages returns x number of public messages that have not been flagged, in desc order by publish date
// "user" is a reserved word in postgres, so it needs to be quoted in the queries
func GetPublicMessages(numberOfMessages int) []models.MessageDto {
	var results []models.MessageDto
	GetPublicMessagesErr := GetDBInstance().db.Model(&models.Message{}).
		Select("message.text, message.pub_date, \"user\".username, \"user\".email").
		Joins("left join \"user\" on message.author_id = \"user\".user_id").
		Where("message.flagged=0").
		Order("pub_date desc").
		Limit(numberOfMessages).
		Scan(&results).Error

	if GetPublicMessagesErr != nil {
		log.WithFields(log.Fields{
			"GetPublicMessagesErr": GetPublicMessagesErr,
			"numberOfMessages":     numberOfMessages,
		}).Error("GetPublicMessages: DB err")
	}

	return results
}

// GetMessagesForUser returns x number of messages for the specified user that have not been flagged,
// in desc order by publish date
// "user" is a reserved word in postgres, so it needs to be quoted in the queries
func GetMessagesForUser(numberOfMessages int, userID uint) []models.MessageDto {
	var results []models.MessageDto
	er := GetDBInstance().db.Model(models.Message{}).
		Order("pub_date desc").
		Select("message.text,message.pub_date, \"user\".email, \"user\".username").
		Joins("left join \"user\" on \"user\".user_id = message.author_id").
		Where("\"user\".user_id=?", userID).
		Limit(numberOfMessages).
		Scan(&results).Error

	if er != nil {
		log.WithFields(log.Fields{
			"err":              er,
			"message":          userID,
			"numberOfMessages": numberOfMessages,
		}).Error("GetMessagesForUser: DB err")
	}

	return results
}

// CreateMessage creates a message in the database
func CreateMessage(message models.Message) error {
	err := GetDBInstance().db.Create(&message).Error
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"message": message,
		}).Error("CreateMessage: DB err")
	}
	return err
}

// GetMessageCount returns the number of messages in the system
func GetMessageCount() int64 {
	var count int64
	GetMessageCountErr := GetDBInstance().db.Find(&models.Message{}).Count(&count).Error
	if GetMessageCountErr != nil {
		log.WithFields(log.Fields{
			"GetMessageCountErr": GetMessageCountErr,
		}).Error("GetMessageCount: DB err")
	}
	return count
}
