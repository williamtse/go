package email

import (
	"gopkg.in/gomail.v2"
	"html/template"
	"strings"
)

type Client struct {
	fromAddress string
	fromName    string
	d           *gomail.Dialer
}

type Conf struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

func NewClient(conf Conf) *Client {
	d := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)
	return &Client{
		fromAddress: conf.Username,
		fromName:    conf.Name,
		d:           d,
	}
}

func (e *Client) RenderTemplate(tmplStr string, pageVariables any) (string, error) {
	// 创建模板并解析 HTML 字符串
	tmpl, err := template.New("verify_code").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var output strings.Builder
	// 执行模板并将结果写入 strings.Builder
	err = tmpl.Execute(&output, pageVariables)
	if err != nil {
		return "", err
	}
	renderedTemplate := output.String()
	// 输出结果
	return renderedTemplate, nil
}

func (e *Client) Send(toAddress string, subject string, tmplStr string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.fromAddress, e.fromName))
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", tmplStr)
	return e.d.DialAndSend(m)
}
