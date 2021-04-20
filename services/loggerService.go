package services

import (
	"runtime/debug"

	"github.com/heyjoakim/devops-21/models"
	log "github.com/sirupsen/logrus"
)

func LogInfo(message string) {
	log.Info(message)
}

func LogPanic(message string){
	log.Panic(message)
}

func LogWarn() {
	log.Warn("warn")
}



func LogError(logObj models.Log) {
	stacktrace := string(debug.Stack())
	logFields := log.Fields{
		"message":        logObj.Message,
		"data":           logObj.Data,
		"additionalInfo": logObj.AdditionalInfo,
		"stacktrace":     stacktrace,
	}
	log.WithFields(logFields).Error(logObj.Message)
}

func LogFatal(logObj models.Log){
	stacktrace := string(debug.Stack())
	logFields := log.Fields{
		"message":        logObj.Message,
		"data":           logObj.Data,
		"additionalInfo": logObj.AdditionalInfo,
		"stacktrace":     stacktrace,
	}
	log.WithFields(logFields).Fatal(logObj.Message)
}

// func executeLog(logFunc func(...interface{}), message models.Log) {
// 	logFunc(message)
// }
