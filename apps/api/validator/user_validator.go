package validator

import (
	"fmt"
	"regexp"
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

type CreateUserInput struct {
	Name  string
	Email string
}

type UpdateUserInput struct {
	Name  string
	Email string
}

// ValidateCreateUser validates input for creating a user
func (v *UserValidator) ValidateCreateUser(input CreateUserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(input.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}

	if len(input.Name) > 100 {
		return fmt.Errorf("name must be less than 100 characters")
	}

	if input.Email == "" {
		return fmt.Errorf("email is required")
	}

	if !isValidEmail(input.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateUpdateUser validates input for updating a user
func (v *UserValidator) ValidateUpdateUser(input UpdateUserInput) error {
	if input.Name != "" {
		if len(input.Name) < 2 {
			return fmt.Errorf("name must be at least 2 characters")
		}

		if len(input.Name) > 100 {
			return fmt.Errorf("name must be less than 100 characters")
		}
	}

	if input.Email != "" && !isValidEmail(input.Email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// isValidEmail checks if email format is valid
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
