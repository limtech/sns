package sns

import (
	"errors"
	"fmt"
	"log"

	"github.com/levigross/grequests"
)

// OAuthWechat oauth wechat
type OAuthWechat struct {
	ClientID     string
	ClientSecret string
}

// NewOAuthWechat return *OAuthWechat
func NewOAuthWechat(id, secret string) *OAuthWechat {
	return &OAuthWechat{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// Authorize for wechat
func (o *OAuthWechat) Authorize(state, callback string) string {
	return fmt.Sprintf(
		`https://open.weixin.qq.com/connect/qrconnect?appid=%s&response_type=code&scope=snsapi_login&state=%s&redirect_uri=%s`,
		o.ClientID,
		state,
		callback,
	)
}

// AccessToken for wechat to get access token
func (o *OAuthWechat) AccessToken(code, state, callback string) (OAuthAccessToken, error) {
	// return data
	rtn := OAuthAccessToken{}
	// response
	resp, err := grequests.Post(
		`https://api.weixin.qq.com/sns/oauth2/access_token`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"appid":      o.ClientID,
				"secret":     o.ClientSecret,
				"code":       code,
				"grant_type": "authorization_code",
			},
		},
	)

	// request faild
	if err != nil {
		log.Println(err)
		return rtn, err
	}

	log.Println(resp.String())

	respData := struct {
		ErrorCode int64  `json:"errcode,omitempty"`
		ErrMsg    string `json:"errmsg,omitempty"`
		// success
		AccessToken  string `json:"access_token,omitempty"`
		ExpiresIn    int64  `json:"expires_in,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
		OpenID       string `json:"openid,omitempty"`
		Scope        string `json:"scope,omitempty"`
		UnionID      string `json:"unionid,omitempty"`
	}{}

	// parse json
	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}

	// faild
	if respData.ErrorCode > 0 {
		return rtn, errors.New(respData.ErrMsg)
	}

	rtn.AccessToken = respData.AccessToken
	rtn.OpenID = respData.OpenID
	rtn.UnionID = respData.UnionID
	return rtn, nil
}

// Userinfo for wechat to get access token
func (o *OAuthWechat) Userinfo(accessToken, openID string) (OAuthUserinfo, error) {
	// return data
	rtn := OAuthUserinfo{
		Platform: OAuthUserPlatformWechat,
		Gender:   OAuthUserGenderUnknown,
	}
	// response
	resp, err := grequests.Get(
		`https://api.weixin.qq.com/sns/userinfo`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"access_token": accessToken,
				"openid":       openID,
				"lang":         "zh_CN",
			},
		},
	)

	if err != nil {
		log.Println(err)
		return rtn, err
	}

	// detail
	rtn.Detail = resp.String()
	log.Println(rtn.Detail)

	// response data
	respData := struct {
		ErrorCode  int64    `json:"errorcode,omitempty"`
		ErrorMsg   string   `json:"errormsg,omitempty"`
		OpenID     string   `json:"openid,omitempty"`
		Nickname   string   `json:"nickname,omitempty"`
		Sex        int      `json:"sex,omitempty"`
		Province   string   `json:"province,omitempty"`
		City       string   `json:"city,omitempty"`
		Country    string   `json:"country,omitempty"`
		HeadImgURL string   `json:"headimgurl,omitempty"`
		Privilege  []string `json:"privilege,omitempty"`
		UnionID    string   `json:"unionid,omitempty"`
	}{}

	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}

	// faild
	if respData.ErrorCode != 0 {
		log.Println(respData.ErrorCode)
		return rtn, errors.New(respData.ErrorMsg)
	}

	// return data
	rtn.OpenID = respData.OpenID
	rtn.UnionID = respData.UnionID
	rtn.Name = respData.Nickname
	rtn.Avatar = respData.HeadImgURL
	// format gender
	switch respData.Sex {
	case 2:
		rtn.Gender = OAuthUserGenderFemale
		break
	case 1:
		rtn.Gender = OAuthUserGenderMale
		break
	default:
		rtn.Gender = OAuthUserGenderUnknown
		break
	}

	return rtn, err
}
