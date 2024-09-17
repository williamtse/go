package douyin

import (
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	douyinGo "github.com/zhangshuai/douyin-go"
	"time"
)

type Client struct {
	manager *douyinGo.Manager
	log     *log.Helper
	conf    Conf
}

type AccessToken struct {
	OpenID                string    `json:"openid"`
	AccessToken           string    `json:"access_token"`
	ExpiredAt             time.Time `json:"expired_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

type Conf struct {
	ClientKey        string
	ClientSecret     string
	DirectURL        string
	ProfileExpiresIn int64
	Scopes           string
}

func NewClient(conf Conf, logger log.Logger) *Client {
	credentials := douyinGo.NewCredentials(conf.ClientKey, conf.ClientSecret)
	return &Client{
		manager: douyinGo.NewManager(credentials, nil),
		log:     log.NewHelper(log.With(logger, "module", "merchant/data")),
		conf:    conf,
	}
}

func (d *Client) GetClientToken() (string, uint64, error) {
	d.log.Infof("生成client token...")
	token, err := d.manager.OauthClientAccessToken()
	if err != nil {
		d.log.Errorf("生成失败 %v", err)
		return "", 0, err
	}
	if token.Message != "success" {
		d.log.Errorf("生成失败 %v", token.Message)
		return "", 0, errors.New(token.Message)
	}
	d.log.Errorf("生成client token 成功 %v", token.Data)
	return token.Data.AccessToken, token.Data.ExpiresIn, nil
}

func (d *Client) GetUserProfile(clientToken string, openId string) (string, time.Duration, error) {
	now := time.Now().Unix()
	secondsToAdd := d.conf.ProfileExpiresIn // 例如，添加 3600 秒（即 1 小时）
	expireTime := now + secondsToAdd
	rs, err := d.manager.SchemaGetUserProfile(douyinGo.SchemaGetUserProfileReq{
		AccessToken: clientToken,
		Body: douyinGo.SchemaGetUserProfileBody{
			ExpireAt: expireTime,
			OpenId:   openId,
		},
	})
	if err != nil {
		d.log.Errorf("接口异常： %v", err)
		return "", 0, err
	}

	if rs.ErrNo > 0 {
		d.log.Errorf("接口报错，错误码： %d， 错误信息：%s", rs.ErrNo, rs.ErrMsg)
		return "", 0, errors.New(rs.ErrMsg)
	}

	d.log.Info("接口返回：", rs)
	return rs.Data.Schema, time.Second * time.Duration(secondsToAdd), nil
}

func (d *Client) RefreshToken(refreshToken string) (*AccessToken, error) {
	accessToken, err := d.manager.OauthRefreshToken(douyinGo.OauthRefreshTokenReq{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	now := time.Now()
	expirationTime := now.Add(time.Duration(accessToken.Data.ExpiresIn) * time.Second)
	refreshTokenExpirationTime := now.Add(time.Duration(accessToken.Data.RefreshExpiresIn) * time.Second)
	return &AccessToken{
		AccessToken:           accessToken.Data.AccessToken,
		ExpiredAt:             expirationTime,
		RefreshTokenExpiredAt: refreshTokenExpirationTime,
		RefreshToken:          accessToken.Data.RefreshToken,
		OpenID:                accessToken.Data.OpenId,
	}, nil
}

func (d *Client) GetAuthUrl() string {
	return d.manager.OauthConnect(douyinGo.OauthParam{
		Scope:       d.conf.Scopes,
		RedirectUri: d.conf.DirectURL,
	})
}

func (d *Client) GetAccessToken(code string) (*AccessToken, error) {
	accessToken, err := d.manager.OauthAccessToken(douyinGo.OauthAccessTokenReq{
		Code: code,
	})
	d.log.Info("获取抖音accessToken接口返回：", accessToken, err)
	if err != nil {
		return nil, fmt.Errorf("获取accessToken失败1:%v", err)
	}
	if accessToken.Message != "success" {
		return nil, fmt.Errorf("获取accessToken失败2:%v", accessToken.Data.Error())
	}
	d.log.Info("获取抖音accessToken成功：", accessToken)
	now := time.Now()
	expirationTime := now.Add(time.Duration(accessToken.Data.ExpiresIn) * time.Second)
	refreshTokenExpirationTime := now.Add(time.Duration(accessToken.Data.RefreshExpiresIn) * time.Second)
	return &AccessToken{
		AccessToken:           accessToken.Data.AccessToken,
		ExpiredAt:             expirationTime,
		RefreshTokenExpiredAt: refreshTokenExpirationTime,
		RefreshToken:          accessToken.Data.RefreshToken,
		OpenID:                accessToken.Data.OpenId,
	}, nil
}

func (d *Client) GetUserInfo(openId string, accessToken string) (string, string, error) {
	info, err := d.manager.OauthUserinfo(douyinGo.OauthUserinfoReq{
		OpenId:      openId,
		AccessToken: accessToken,
	})

	if err != nil {
		return "", "", err
	}
	return info.Data.Nickname, info.Data.Avatar, nil
}

func (d *Client) GetUserFans(openId string, accessToken string, days int64) (int64, error) {
	opts := douyinGo.DataExternalUserFansReq{
		OpenId:      openId,
		AccessToken: accessToken,
		DataType:    days,
	}
	d.log.Infof("请求参数： DataType: %d, OpenId: %s, AccessToken: %s", days, openId, accessToken)
	rs, err := d.manager.DataExternalUserFans(opts)
	if err != nil {
		d.log.Infof("获取用户粉丝接口异常:%v", err)
		return 0, err
	}
	errMsg := rs.Data.Error()
	if rs.Data.ErrorCode != 0 {
		d.log.Infof("获取用户粉丝接口返回错误:%s, %s", errMsg, rs.Extra.SubDescrition)
		return 0, err
	}
	if len(rs.Data.ResultList) == 0 {
		return 0, err
	}
	total := rs.Data.ResultList[0].TotalFans
	return total, err
}
