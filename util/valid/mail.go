package valid

import "regexp"

func IsValidEmail(email string) bool {
	// 正则表达式来验证邮箱地址格式
	// 此处使用较为简单的示例，实际使用中可以根据需求调整正则表达式
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
