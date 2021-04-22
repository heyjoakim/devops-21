package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/heyjoakim/devops-21/models"
	"golang.org/x/crypto/bcrypt"
)

// GetUserID returns user ID for username
func GetUserID(username string) (uint, error) {
	var user models.User
	getUserIDErr := GetDBInstance().db.First(&user, "username = ?", username).Error
	if getUserIDErr != nil {
		logUserErr(getUserIDErr, username)
	}
	return user.UserID, getUserIDErr
}

func GetUserFromUsername(username string) (models.User, error) {
	var user models.User
	err := GetDBInstance().db.Where("username = ?", username).First(&user).Error
	if err != nil {
		logUserErr(err, username)
	}
	return user, err
}

func GetUser(userID uint) models.User {
	var user models.User
	getUserErr := GetDBInstance().db.First(&user, "user_id = ?", userID).Error
	if getUserErr != nil {
		logUserErr(getUserErr, userID)
	}
	return user
}

// CreateUser creates a new user in the database
func CreateUser(userRequest models.UserCreateRequest) error {
	err := validateUserCreateRequest(userRequest)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		logUserErr(err, userRequest)
		return err
	}
	user := models.User{
		Username: userRequest.Username,
		Email:    userRequest.Email,
		PwHash:   string(hash),
	}
	createUserErr := GetDBInstance().db.Create(&user).Error
	if createUserErr != nil {
		logUserErr(createUserErr, user)
	}
	return createUserErr
}

// GetUserCount returns the number of users reigstered in the system
func GetUserCount() int64 {
	var count int64
	countErr := GetDBInstance().db.Find(&models.User{}).Count(&count).Error
	if countErr != nil {
		LogError(models.Log{
			Message: countErr.Error(),
		})
	}
	return count
}

func LoginUser(loginRequest models.LoginRequest) (models.User, error) {
	user, err := GetUserFromUsername(loginRequest.Username)
	if err != nil {
		LogWarn(fmt.Sprintf("Login attempt with unknown user: %s", loginRequest.Username))
		return models.User{}, errors.New("error in username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PwHash), []byte(loginRequest.Password)); err != nil {
		LogWarn(fmt.Sprintf("Login attempt with wrong password for user: %s", loginRequest.Username))
		return models.User{}, errors.New("error in username or password")
	}

	return user, nil
}

func logUserErr(err error, data interface{}) {
	LogError(models.Log{
		Message: err.Error(),
		Data:    data,
	})
}

func validateUserCreateRequest(user models.UserCreateRequest) error {
	if err := validateUsername(user.Username); err != nil {
		return err
	}
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	if err := validatePassword(user.Password, user.Password2); err != nil {
		return err
	}
	return nil
}

func validateUsername(username string) error {
	const charLimit = 20
	if len(username) == 0 {
		return errors.New("you have to enter a username")
	}
	if len(username) > charLimit {
		return errors.New("username cannot be longer than 20 characters")
	}
	if _, err := GetUserID(username); err == nil {
		return errors.New("the username is already taken")
	}
	return nil
}

func validatePassword(password string, password2 string) error {
	const charLimit = 20
	if len(password) > charLimit {
		return errors.New("password cannot be longer than 20 characters")
	}
	if len(password) == 0 {
		return errors.New("you have to enter a password")
	}
	if password != password2 {
		return errors.New("the two passwords do not match")
	}
	return nil
}

func validateEmail(email string) error {
	const emailLimit = 30
	if len(email) == 0 || !strings.Contains(email, "@") || len(email) > emailLimit {
		return errors.New("you have to enter a valid email address")
	}
	return nil
}
