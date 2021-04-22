package services

import (
	"runtime/debug"

	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

func LogInfo(message string) {
	log.Info(message)
}

func LogPanic(message string) {
	log.Panic(message)
}

func LogWarn() {
	log.Warn("warn")
}

func LogError(logObj models.Log) {
	logFields := makeLogFields(logObj)
	log.WithFields(logFields).Error(logObj.Message)
}

func LogFatal(logObj models.Log) {
	logFields := makeLogFields(logObj)
	log.WithFields(logFields).Fatal(logObj.Message)
}

func makeLogFields(msg models.Log) log.Fields {
	stacktrace := string(debug.Stack())
	logFields := log.Fields{
		"message":        msg.Message,
		"data":           msg.Data,
		"additionalInfo": msg.AdditionalInfo,
		"stacktrace":     stacktrace,
	}
	return logFields
}
