package email

import (
	"gopkg.in/gomail.v2"
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

func (e *Client) Send(toAddress string, subject string, tmplStr string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(e.fromAddress, e.fromName))
	m.SetHeader("To", toAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", tmplStr)
	return e.d.DialAndSend(m)
}
