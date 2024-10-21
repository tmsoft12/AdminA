package service

import (
	"errors"
	"rr/domain"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtSecret = []byte("your_secret_key")

// Ulanyjyny bcrypt bilen şifrelenen paroly bilen hasaba al
func RegisterUser(db *gorm.DB, username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := domain.User{Username: username, Password: hashedPassword}
	return db.Create(&user).Error
}

// Ulanyjyny tassyklap, JWT token döret
func LoginUser(db *gorm.DB, username string, password string) (string, error) {
	var user domain.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return "", errors.New("nädogry ulanyjy ýa-da parol")
	}

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", errors.New("nädogry ulanyjy ýa-da parol")
	}

	// JWT token döret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 sagat möhlet
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// JWT token-y tassyklamak
func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("nädogry ýazuw usuly")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
