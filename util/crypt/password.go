package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

func HashCheck(password string, hashedPassword string) bool {
	// 比较哈希密码和用户输入的密码
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}

func HashMake(password string) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPassword
}
