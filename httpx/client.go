package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPClient 封装 HTTP 请求
type HTTPClient struct {
	Client  *http.Client
	Headers map[string]string
}

// NewHTTPClient 创建一个新的 HTTPClient 实例
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client:  &http.Client{},
		Headers: make(map[string]string),
	}
}

// SetHeader 设置请求头
func (c *HTTPClient) SetHeader(key, value string) {
	c.Headers[key] = value
}

// Do 发送 HTTP 请求
func (c *HTTPClient) Do(method, url string, body interface{}) (*http.Response, error) {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	return c.Client.Do(req)
}
