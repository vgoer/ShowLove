// Package validator provides common validation utilities for user inputs.
package validator

import (
	"errors"
	"net/mail"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidEmail    = errors.New("邮箱格式不正确")
	ErrInvalidPassword = errors.New("密码需要至少8位，包含字母和数字")
	ErrInvalidNickname = errors.New("昵称需要1-100个字符且不能全为空格")
	ErrEmptyContent    = errors.New("内容不能为空")
	ErrContentTooLong  = errors.New("内容超出字数限制")
)

// ValidateEmail checks if the email has a valid format.
func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword checks if the password meets minimum requirements.
// Must be at least 8 characters, contain both letters and numbers.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	hasLetter := false
	hasDigit := false
	for _, c := range password {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return ErrInvalidPassword
	}
	return nil
}

// ValidateNickname checks if the nickname is valid (1-100 chars, not all whitespace).
func ValidateNickname(nickname string) error {
	trimmed := strings.TrimSpace(nickname)
	if trimmed == "" {
		return ErrInvalidNickname
	}
	if utf8.RuneCountInString(nickname) > 100 {
		return ErrInvalidNickname
	}
	return nil
}

// ValidateContent checks if the content is valid (non-empty, within maxLen).
func ValidateContent(content string, maxLen int) error {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ErrEmptyContent
	}
	if utf8.RuneCountInString(content) > maxLen {
		return ErrContentTooLong
	}
	return nil
}
