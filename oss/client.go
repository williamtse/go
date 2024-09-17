package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/williamtse/gopkg/encrypt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Client 提供图片上传功能的 SDK
type Client struct {
	host      string
	uploadUrl string // 上传接口的 URL
	keyID     string
	keySecret string
}

// NewClient 创建一个新的 Client 实例
func NewClient(host, keyId, keySecret string) *Client {
	return &Client{
		host:      host,
		uploadUrl: host + "/upload",
		keyID:     keyId,
		keySecret: keySecret,
	}
}

func (client *Client) GetHost() string {
	return client.host
}

// Upload 上传图片到指定的接口
func (iu *Client) Upload(imageURL string) (string, error) {
	// 下载图片并保存到临时文件
	tempFile, err := downloadImage(imageURL)
	if err != nil {
		return "", fmt.Errorf("error downloading image: %v", err)
	}
	defer os.Remove(tempFile.Name()) // 删除临时文件

	// 准备上传图片到接口
	url, err := iu.uploadFile(iu.uploadUrl, tempFile)
	if err != nil {
		return "", fmt.Errorf("error uploading image: %v", err)
	}

	fmt.Println("Image uploaded successfully.")
	return url, nil
}

// downloadImage 下载图片并保存到临时文件，返回临时文件的指针和可能的错误
func downloadImage(url string) (*os.File, error) {
	// 发起 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 创建临时文件来保存下载的图片
	file, err := ioutil.TempFile("", "download-*.jpg")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 将 HTTP 响应的 Body 写入临时文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Image downloaded:", file.Name())
	return file, nil
}

func ExtractBase64Data(dataURL string) (string, error) {
	// 找到 Base64 数据的开始部分
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid data URL format")
	}

	// 返回 Base64 编码数据
	return parts[1], nil
}

// 转发文件到另一个服务器
func (iu *Client) UploadFileByImageData(filename string, base64Url string) (string, error) {
	base64Data, err := ExtractBase64Data(base64Url)
	if err != nil {
		return "", err
	}
	// 解码 Base64 字符串
	imgData, err := encrypt.Base64Decode(base64Data)
	if err != nil {
		return "", err
	}
	// 创建新的 multipart/form-data 请求
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	fileWriter, err := writer.CreateFormFile("image", filename) // 这里设置文件名
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// 将解码后的图像数据写入文件字段
	_, err = fileWriter.Write(imgData)
	if err != nil {
		return "", fmt.Errorf("failed to write image data: %w", err)
	}

	// 关闭 multipart writer，结束请求的组装
	err = writer.Close()
	if err != nil {
		return "", err
	}

	return iu.uploadFileByBuffer(iu.uploadUrl, writer, body)
}

// uploadFile 上传文件到目标接口
func (c *Client) uploadFile(uploadURL string, file *os.File) (string, error) {
	// 打开要上传的文件
	fileContents, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	// 创建一个 buffer 用于组装 multipart/form-data 请求
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// 添加文件字段到 multipart 请求
	part, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}
	_, err = part.Write(fileContents)
	if err != nil {
		return "", err
	}

	// 关闭 multipart writer，结束请求的组装
	err = writer.Close()
	if err != nil {
		return "", err
	}

	// 发起 HTTP POST 请求上传文件
	return c.uploadFileByBuffer(uploadURL, writer, body)
}

func (c *Client) uploadFileByBuffer(uploadURL string, writer *multipart.Writer, body *bytes.Buffer) (string, error) {
	// 发起 HTTP POST 请求上传文件
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	timestamp := time.Now().Format(time.RFC3339)
	stringToSign := fmt.Sprintf("%s\n%s\n%s", "POST", "/upload", "timestamp="+timestamp)
	signature := calculateSignature(stringToSign, c.keySecret)
	auth := fmt.Sprintf("accessKeyID=%s,timestamp=%s,signature=%s", c.keyID, timestamp, signature)
	// 添加身份验证信息到请求头部
	req.Header.Set("Authorization", auth)
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 检查上传是否成功
	if resp.StatusCode != http.StatusOK {

		return "", fmt.Errorf("upload failed: %s", resp.Status)
	}

	resbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResp map[string]interface{}
	err = json.Unmarshal(resbody, &apiResp)
	if err != nil {
		return "", err
	}
	if url, ok := apiResp["data"].(string); ok {
		return c.host + url, nil
	}
	return "", fmt.Errorf("response data error")
}

// 计算签名
func calculateSignature(stringToSign, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
