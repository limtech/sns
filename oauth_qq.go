package sns

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/levigross/grequests"
)

// OAuthQQ oauth qq
type OAuthQQ struct {
	ClientID     string
	ClientSecret string
}

// NewOAuthQQ return *OAuthQQ
func NewOAuthQQ(id, secret string) *OAuthQQ {
	return &OAuthQQ{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// Authorize for qq
func (o *OAuthQQ) Authorize(state, callback string) string {
	return fmt.Sprintf(
		`https://graph.qq.com/oauth2.0/authorize?client_id=%s&response_type=code&state=%s&redirect_uri=%s`,
		o.ClientID,
		state,
		callback,
	)
}

// AccessToken for qq to get access token
func (o *OAuthQQ) AccessToken(code, state, callback string) (OAuthAccessToken, error) {
	// return data
	rtn := OAuthAccessToken{}
	// request
	resp, err := grequests.Get(
		`https://graph.qq.com/oauth2.0/token`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     o.ClientID,
				"client_secret": o.ClientSecret,
				"code":          code,
				"redirect_uri":  callback,
			},
		},
	)
	// request faild
	if err != nil {
		log.Println(err)
		return rtn, err
	}

	// error jsonp to json
	respJSON := regexp.MustCompile(`callback\((.*)\)(;)?`).ReplaceAllString(resp.String(), "$1")
	respJSONData := struct {
		Error            int    `json:"error"`
		ErrorDescription string `json:"error_description"`
	}{}

	// error
	if json.Unmarshal([]byte(respJSON), &respJSONData) == nil {
		log.Println(respJSONData)
		return rtn, errors.New(respJSONData.ErrorDescription)
	}

	// parse url query
	respData, err := url.ParseQuery(resp.String())
	if err != nil {
		return rtn, err
	}

	rtn.AccessToken = respData.Get("access_token")
	rtn.OpenID, err = o.GetOpenID(rtn.AccessToken) // get openid
	if err != nil {
		return rtn, err
	}
	return rtn, nil
}

// GetOpenID return openid by accesstoken
func (o *OAuthQQ) GetOpenID(accessToken string) (string, error) {
	openid := ""
	// request
	resp, err := grequests.Get(
		`https://graph.qq.com/oauth2.0/me`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"access_token": accessToken,
			},
		},
	)

	// request faild
	if err != nil {
		log.Println(err)
		return openid, err
	}

	// jsonp to json
	respJSON := regexp.MustCompile(`callback\((.*)\)(;)?`).ReplaceAllString(resp.String(), "$1")
	// responseJson
	respData := struct {
		// error
		Error            int    `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		// success
		ClientID string `json:"client_id"`
		OpenID   string `json:"openid"`
	}{}

	// request faild
	if err := json.Unmarshal([]byte(respJSON), &respData); err != nil {
		log.Println(err)
		return openid, err
	}

	// error
	if respData.Error > 0 {
		log.Println(err)
		return openid, errors.New(respData.ErrorDescription)
	}

	return respData.OpenID, nil
}

// Userinfo for qq by access token
func (o *OAuthQQ) Userinfo(accessToken, openID string) (OAuthUserinfo, error) {
	// return data
	rtn := OAuthUserinfo{
		Platform: OAuthUserPlatformQQ,
		Gender:   OAuthUserGenderUnknown,
	}

	// request
	resp, err := grequests.Get(
		`https://graph.qq.com/user/get_user_info`,
		&grequests.RequestOptions{
			Params: map[string]string{
				"access_token":       accessToken,
				"oauth_consumer_key": o.ClientID,
				"openid":             openID,
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

	// responseJson
	respData := struct {
		Ret int    `json:"ret,omitempty"`
		Msg string `json:"Met,omitempty"`
		// success
		ID       string `json:"id"`
		Avatar   string `json:"figureurl_qq"`
		AvatarHD string `json:"figureurl_qq_2"`
		Nickname string `json:"nickname"`
		Gender   string `json:"gender"`
	}{}

	// request faild
	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}

	// faild
	if respData.Ret != 0 {
		return rtn, errors.New(respData.Msg)
	}

	// return data
	rtn.OpenID = openID
	rtn.Name = respData.Nickname
	rtn.Avatar = respData.Avatar
	// format gender
	switch respData.Gender {
	case "女":
		rtn.Gender = OAuthUserGenderFemale
		break
	case "男":
		rtn.Gender = OAuthUserGenderMale
		break
	default:
		rtn.Gender = OAuthUserGenderUnknown
		break
	}

	return rtn, nil
}
