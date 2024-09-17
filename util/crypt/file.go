package crypt

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetFilenameBase64MD5(base64Str string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(base64Str))
	if err != nil {
		return "", err
	}

	fileFormat, err := ExtractFileFormat(base64Str)
	if err != nil {
		return "", err
	}
	md5str := hex.EncodeToString(hash.Sum(nil))
	newFileName := fmt.Sprintf("%s.%s", md5str, fileFormat)
	return newFileName, nil
}

func ExtractFileFormat(base64Str string) (string, error) {
	// 检查 Base64 字符串是否包含数据 URI 前缀
	if !strings.HasPrefix(base64Str, "data:") {
		return "", fmt.Errorf("无效的 Base64 字符串: 缺少数据 URI 前缀")
	}

	// 找到 MIME 类型的结束位置
	commaIndex := strings.Index(base64Str, ",")
	if commaIndex == -1 {
		return "", fmt.Errorf("无效的 Base64 字符串: MIME 类型和数据部分之间缺少逗号")
	}

	// 提取 MIME 类型部分
	mimeTypePart := base64Str[5:commaIndex]
	mimeTypeParts := strings.Split(mimeTypePart, ";")

	// MIME 类型的第一部分通常是 "image/*" 格式
	if len(mimeTypeParts) == 0 {
		return "", fmt.Errorf("无法提取 MIME 类型")
	}

	// 提取文件格式
	mimeType := mimeTypeParts[0]
	fileFormat := strings.TrimPrefix(mimeType, "image/")
	if fileFormat == mimeType {
		return "", fmt.Errorf("无法提取文件格式")
	}

	return fileFormat, nil
}
