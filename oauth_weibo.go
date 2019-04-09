package sns

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/levigross/grequests"
)

// OAuthWeibo oauth weibo
type OAuthWeibo struct {
	ClientID     string
	ClientSecret string
}

// NewOAuthWeibo return *OAuthWeibo
func NewOAuthWeibo(id, secret string) *OAuthWeibo {
	return &OAuthWeibo{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// Authorize for weibo
func (o *OAuthWeibo) Authorize(state, callback string) string {
	return fmt.Sprintf(
		`https://api.weibo.com/oauth2/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s`,
		o.ClientID,
		state,
		callback,
	)
}

// AccessToken for weibo to get access token
func (o *OAuthWeibo) AccessToken(code, state, callback string) (OAuthAccessToken, error) {
	// return data
	rtn := OAuthAccessToken{}
	// response
	resp, err := grequests.Post(
		`https://api.weibo.com/oauth2/access_token`,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Accept": "application/json",
			},
			Data: map[string]string{
				"client_id":     o.ClientID,
				"client_secret": o.ClientSecret,
				"code":          code,
				"redirect_uri":  callback,
				"grant_type":    "authorization_code",
			},
		},
	)
	// request faild
	if err != nil {
		log.Println(err)
		return rtn, err
	}

	log.Println(resp.String())

	// return data
	respData := struct {
		Error            string `json:"error,omitempty"`
		ErrorCode        int64  `json:"error_code,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		// success
		AccessToken string `json:"access_token"`
		RemindIn    string `json:"remind_in"`
		ExpiresIn   int64  `json:"expires_in"`
		UID         string `json:"uid"`
		IsRealName  string `json:"isRealName"`
	}{}
	// return json
	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}
	// error
	if respData.ErrorCode > 0 {
		return rtn, errors.New(respData.ErrorDescription)
	}

	rtn.AccessToken = respData.AccessToken
	rtn.OpenID = respData.UID
	return rtn, nil
}

// Userinfo for weibo to get access token
func (o *OAuthWeibo) Userinfo(accessToken, openID string) (OAuthUserinfo, error) {
	// return data
	rtn := OAuthUserinfo{
		Platform: OAuthUserPlatformWeibo,
		Gender:   OAuthUserGenderUnknown,
	}
	// response
	resp, err := grequests.Get(
		`https://api.weibo.com/2/users/show.json`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"access_token": accessToken,
				"uid":          openID,
			},
		},
	)
	// request faild
	if err != nil {
		log.Println(err)
		return rtn, err
	}

	// detail
	rtn.Detail = resp.String()
	log.Println(rtn.Detail)

	respData := struct {
		Error     string `json:"error"`
		ErrorCode int64  `json:"error_code"`
		Request   string `json:"request"`
		// success
		ID     int64  `json:"id"`
		Name   string `json:"screen_name"`
		Avatar string `json:"avatar_hd"`
		Gender string `json:"gender"`
	}{}

	// parse error
	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}
	// auth faild
	if respData.ErrorCode > 0 {
		return rtn, errors.New(respData.Error)
	}

	// return data
	rtn.OpenID = strconv.FormatInt(respData.ID, 10)
	rtn.Name = respData.Name
	rtn.Avatar = respData.Avatar
	// format gender
	switch respData.Gender {
	case "f":
		rtn.Gender = OAuthUserGenderFemale
		break
	case "m":
		rtn.Gender = OAuthUserGenderMale
		break
	default:
		rtn.Gender = OAuthUserGenderUnknown
		break
	}

	return rtn, err
}
