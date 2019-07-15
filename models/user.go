package models

import (
	"errors"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"unicode"
)

type User struct {
	Id              int    `json:"id"`
	BusinessName    string `json:"business_name"`
	FullName        string `json:"full_name"`
	BusinessEmail   string `json:"business_email"`
	BusinessPhone   string `json:"business_phone"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password,omitempty"gorm:"-"`
	IsVerified      bool   `json:"-"`
	EmailCode       string `json:"-"`
	PasswordCode    string `json:"-"`
}

func (User) TableName() string {
	return "tbl_users"
}

func (user *User) Validate() error {
	return validation.ValidateStruct(user,
		validation.Field(&user.BusinessName, validation.Required),
		validation.Field(&user.FullName, validation.Required),
		validation.Field(&user.BusinessEmail, validation.Required, is.Email),
		validation.Field(&user.BusinessPhone, validation.Required),
		validation.Field(&user.Password, validation.Required, validation.RuneLength(8, 100), validation.By(checkPassword)),
		validation.Field(&user.ConfirmPassword, validation.Required, validation.By(user.checkConfirmPass)),
	)
}

func checkPassword(value interface{}) error {
	pass, ok := value.(string)
	if !ok {
		return errors.New("check password")
	}

	var (
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
	)
	if len(pass) < 8 {
		return errors.New("password is too short")
	}
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must have at least 1 capital letter")
	}
	if !hasDigit {
		return errors.New("password must have at least 1 digit")
	}
	if !hasSpecial {
		return errors.New("password must have at least 1 special character")
	}
	return nil
}

func (user *User) checkConfirmPass(interface{}) error {
	if user.ConfirmPassword != user.Password {
		return errors.New("password and confirm password does not match")
	}

	return nil
}
