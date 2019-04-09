package sns

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/levigross/grequests"
)

// OAuthGithub oauth github
type OAuthGithub struct {
	ClientID     string
	ClientSecret string
}

// NewOAuthGithub return *OAuthGithub
func NewOAuthGithub(id, secret string) *OAuthGithub {
	return &OAuthGithub{
		ClientID:     id,
		ClientSecret: secret,
	}
}

// Authorize for github
func (o *OAuthGithub) Authorize(state, callback string) string {
	return fmt.Sprintf(
		`https://github.com/login/oauth/authorize?client_id=%s&state=%s&redirect_uri=%s`,
		o.ClientID,
		state,
		callback,
	)
}

// AccessToken for github to get access token
func (o *OAuthGithub) AccessToken(code, state, callback string) (OAuthAccessToken, error) {
	// return data
	rtn := OAuthAccessToken{}
	// request
	resp, err := grequests.Post(
		`https://github.com/login/oauth/access_token`,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Accept": "application/json",
			},
			Data: map[string]string{
				"client_id":     o.ClientID,
				"client_secret": o.ClientSecret,
				"code":          code,
				"redirect_uri":  callback,
				"state":         state,
			},
		},
	)
	// request faild
	if err != nil {
		log.Println(err)
		return rtn, err
	}

	respData := struct {
		Error            string `json:"error,omitempty"`
		ErrorDescription string `json:"error_description,omitempty"`
		ErrorURI         string `json:"error_uri,omitempty"`
		// success
		AccessToken string `json:"access_token,omitempty"`
		Scope       string `json:"scope,omitempty"`
		TokenType   string `json:"token_type,omitempty"`
	}{}

	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}
	if respData.Error != "" {
		log.Println(respData.ErrorDescription)
		return rtn, errors.New(respData.ErrorDescription)
	}

	rtn.AccessToken = respData.AccessToken
	return rtn, nil
}

// Userinfo for github by access token
func (o *OAuthGithub) Userinfo(accessToken, openID string) (OAuthUserinfo, error) {
	// return data
	rtn := OAuthUserinfo{
		Platform: OAuthUserPlatformGithub,
		Gender:   OAuthUserGenderUnknown,
	}
	// request
	resp, err := grequests.Get(
		`https://api.github.com/user`,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": fmt.Sprintf("bearer %s", accessToken),
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
		ID        int64  `json:"id"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
		Message   string `json:"message,omitempty"`
	}{}

	if err := resp.JSON(&respData); err != nil {
		log.Println(err)
		return rtn, err
	}
	if respData.Message != "" {
		log.Println(respData.Message)
		return rtn, errors.New(respData.Message)
	}

	// return data
	rtn.OpenID = strconv.FormatInt(respData.ID, 10)
	rtn.Name = respData.Name
	rtn.Avatar = respData.AvatarURL
	return rtn, nil
}
