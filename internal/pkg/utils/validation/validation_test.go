package validation

import (
	"testing"

	"github.com/K1tten2005/avito_pvz/internal/models"
)

func TestIsValidProductType(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"электроника", true},
		{"обувь", true},
		{"мебель", false},
	}

	for _, tt := range tests {
		if got := IsValidProductType(tt.input); got != tt.want {
			t.Errorf("IsValidProductType(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsValidCity(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Москва", true},
		{"Санкт Петербург", true},
		{"Новосибирск", false},
	}

	for _, tt := range tests {
		if got := IsValidCity(tt.input); got != tt.want {
			t.Errorf("IsValidCity(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{models.RoleEmployee, true},
		{models.RoleModerator, true},
		{"admin", false},
	}

	for _, tt := range tests {
		if got := IsValidRole(tt.input); got != tt.want {
			t.Errorf("IsValidRole(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	salt := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	password := "StrongPass1!"

	hashed := HashPassword(salt, password)

	if !CheckPassword(hashed, password) {
		t.Error("CheckPassword() failed: correct password not recognized")
	}

	if CheckPassword(hashed, "WrongPass123") {
		t.Error("CheckPassword() failed: wrong password recognized as correct")
	}
}

func TestValidEmail(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"user@example.com", true},
		{"invalid-email", false},
		{"user@.com", false},
	}

	for _, tt := range tests {
		if got := ValidEmail(tt.input); got != tt.want {
			t.Errorf("ValidEmail(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}

func TestValidPassword(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Aa1!aaaa", true},
		{"short1!", false},        
		{"alllowercase1!", false}, 
		{"ALLUPPERCASE1!", false}, 
		{"NoSpecialChar1", false}, 
		{"NoDigit!Aa", false},     
	}

	for _, tt := range tests {
		if got := ValidPassword(tt.input); got != tt.want {
			t.Errorf("ValidPassword(%q) = %v; want %v", tt.input, got, tt.want)
		}
	}
}
