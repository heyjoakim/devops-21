package services

import (
	"strconv"

	"github.com/heyjoakim/devops-21/models"
)

func UpdateLatest(latest int) {
	var c models.Config

	er := d.db.First(&c, "key = ?", "latest").Error
	if er != nil {
		c.ID = 0
		c.Key = "latest"
		c.Value = strconv.Itoa(latest)
		d.db.Create(&c)
	} else {
		d.db.Model(&models.Config{}).Where("key = ?", "latest").Update("Value", latest)
	}
}

func GetLatest() int {
	var result int

	d.db.Model(models.Config{}).First(&result, "key = ?", "latest")

	return result
}
