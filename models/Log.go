package models

// Log represents a log message object
type Log struct {
	Message        string
	Data           interface{}
	AdditionalInfo interface{}
}
