package validation

import (
	"bytes"
	"regexp"
	"unicode"

	"github.com/K1tten2005/avito_pvz/internal/models"
	"golang.org/x/crypto/argon2"
)

const (
	maxEmailLength = 20
	minEmailLength = 3
	minPassLength  = 8
	maxPassLength  = 25
)

var AllowedCities = []string{
	"Москва",
	"Санкт Петербург",
	"Казань",
}

func IsValidCity(city string) bool {
	for _, val := range AllowedCities {
		if city == val {
			return true
		}
	}
	return false
}

func IsValidRole(role string) bool {
	return role == models.RoleEmployee || role == models.RoleModerator
}

func HashPassword(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), salt, 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func CheckPassword(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[:8])
	userPassHash := HashPassword(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidPassword(password string) bool {
	var up, low, digit, special bool

	if len(password) < minPassLength || len(password) > maxPassLength {
		return false
	}

	for _, char := range password {

		switch {
		case unicode.IsUpper(char):
			up = true
		case unicode.IsLower(char):
			low = true
		case unicode.IsDigit(char):
			digit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		default:
			return false
		}
	}

	return up && low && digit && special
}
