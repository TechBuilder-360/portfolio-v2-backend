package util

import (
	"errors"
	"fmt"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"
	p "github.com/go-passwd/validator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

func AddrToString(txt *string) string {
	if txt != nil {
		return *txt
	}

	return ""
}

func ToLower(txt string) string {
	return strings.ToLower(txt)
}

func ValidateEmailAddress(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func GenerateUUID() string {
	return uuid.NewString()
}

func ValidatePassword(password string) error {
	passwordValidator := p.New(p.MinLength(8, errors.New("password minimum length is 8 character")))
	err := passwordValidator.Validate(password)
	if err != nil {
		return err
	}

	return nil
}

func HashPassword(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

func msgForTag(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "oneof":
		{
			p := strings.ReplaceAll(param, " ", ",")
			return fmt.Sprintf("This field accepts one of %s", p)
		}
	}
	return ""
}

func CustomErrorResponse(err error) *[]types.ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]types.ApiError, len(ve))
			for i, fe := range ve {
				out[i] = types.ApiError{fe.Field(), msgForTag(fe.Tag(), fe.Param())}
			}
			return &out
		}
	}
	return nil
}

func CapitalizeFirstCharacter(s string) string {
	return cases.Title(language.AmericanEnglish, cases.NoLower).String(strings.ToLower(strings.TrimSpace(s)))
}
