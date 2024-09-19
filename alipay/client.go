package alipay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"github.com/go-pay/xlog"
	"log"
	"net/http"
)

type PaymentParams struct {
	TradeNo    string
	Amount     float32
	TimeExpire string
	Subject    string
}

type Payment interface {
	PagePay(ctx context.Context, params PaymentParams) (string, error)
	ParseNotifyData(r *http.Request) ([]byte, error)
}

type Client struct {
	client *alipay.Client
	Config AlipayConf
}

type AlipayNotifyReq struct {
	NotifyTime    string `json:"notify_time"` // 使用 time.Time 类型
	NotifyType    string `json:"notify_type"`
	NotifyID      string `json:"notify_id"`
	SignType      string `json:"sign_type"`
	Sign          string `json:"sign"`
	TradeNo       string `json:"trade_no"`
	AppID         string `json:"app_id"`
	AuthAppID     string `json:"auth_app_id"`
	OutTradeNo    string `json:"out_trade_no"`
	TradeStatus   string `json:"trade_status"`
	TotalAmount   string `json:"total_amount"`   // 使用 float32 类型以匹配 Protobuf 的 float 类型
	ReceiptAmount string `json:"receipt_amount"` // 使用 float32 类型以匹配 Protobuf 的 float 类型
}

const (
	TradeStatusWaitBuyerPay  = "WAIT_BUYER_PAY" // 交易创建，等待买家付款
	TradeStatusTradeClosed   = "TRADE_CLOSED"   // 未付款交易超时关闭，或支付完成后全额退款
	TradeStatusTradeSuccess  = "TRADE_SUCCESS"  // 交易支付成功
	TradeStatusTradeFinished = "TRADE_FINISHED" // 交易结束，不可退款
)

type AlipayConf struct {
	AppID      string `json:"app_id"`
	PrivateKey string `json:"private_key"`
	ReturnURL  string `json:"return_url"`
	NotifyURL  string `json:"notify_url"`
	AppCert    string `json:"app_cert"`
	RootCert   string `json:"root_cert"`
	PublicCert string `json:"public_cert"`
	ExpiresIn  int64
}

func NewClient(conf AlipayConf) (*Client, error) {
	// 初始化支付宝客户端
	// appid：应用ID
	// privateKey：应用私钥，支持PKCS1和PKCS8
	// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
	log.Println("初始化支付宝客户端")
	confJson, err := json.Marshal(conf)
	log.Println("配置：", string(confJson))
	client, err := alipay.NewClient(conf.AppID, conf.PrivateKey, true)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}

	// 打开Debug开关，输出日志，默认关闭
	client.DebugSwitch = gopay.DebugOn

	// 设置支付宝请求 公共参数
	//    注意：具体设置哪些参数，根据不同的方法而不同，此处列举出所有设置参数
	client.SetLocation(alipay.LocationShanghai). // 设置时区，不设置或出错均为默认服务器时间
							SetCharset(alipay.UTF8).      // 设置字符编码，不设置默认 utf-8
							SetSignType(alipay.RSA2).     // 设置签名类型，不设置默认 RSA2
							SetReturnUrl(conf.ReturnURL). // 设置返回URL
							SetNotifyUrl(conf.NotifyURL)  //. // 设置异步通知URL

	err = client.SetCertSnByPath(conf.AppCert, conf.RootCert, conf.PublicCert)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	return &Client{
		client: client,
		Config: conf,
	}, nil
}

func (c *Client) PagePay(ctx context.Context, args PaymentParams) (string, error) {
	mb := gopay.BodyMap{
		"out_trade_no": args.TradeNo,
		"total_amount": args.Amount,
		"subject":      args.Subject,
		"product_code": "FAST_INSTANT_TRADE_PAY",
		"time_expire":  args.TimeExpire,
	}
	return c.client.TradePagePay(ctx, mb)
}

func (c *Client) ParseNotifyData(r *http.Request) ([]byte, error) {
	notifyReq, err := alipay.ParseNotifyToBodyMap(r)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	payload, err := json.Marshal(notifyReq)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	// 支付宝异步通知验签（公钥证书模式）
	log.Println("回调参数：", string(payload))
	log.Println("签名验证，证书：", c.Config.PublicCert)
	ok, err := alipay.VerifySignWithCert(c.Config.PublicCert, notifyReq)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("验签失败")
	}
	log.Println("签名验证成功！")
	return payload, nil
}
