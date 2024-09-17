package encrypt

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// 生成哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hash string) bool {
	// 比较密码和哈希值
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false // 密码不匹配
	}
	return true // 密码匹配
}
