package services

import (
	"strconv"

	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

func UpdateLatest(latest int) {
	var c models.Config

	err := d.db.First(&c, "key = ?", "latest").Error
	if err != nil {
		log.WithField("err", err).Error("Latest does not exist: DB err")
		c.ID = 0
		c.Key = "latest"
		c.Value = strconv.Itoa(latest)
		d.db.Create(&c)
	} else {
		err := d.db.Model(&models.Config{}).Where("key = ?", "latest").Update("Value", latest).Error
		if err != nil {
			log.WithField("err", err).Error("UpdateLatest: DB err")
		}
	}
}

func GetLatest() int {
	var result int
	err := d.db.Model(models.Config{}).Select("value").First(&result).Error
	if err != nil {
		log.WithField("err", err).Error("GetLatest: DB err")
	}
	return result
}
