package helpers

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

func NewCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 6*time.Second)
}
func IsValidPhoneNumber(phoneNumber string) bool {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	return re.MatchString(phoneNumber)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
