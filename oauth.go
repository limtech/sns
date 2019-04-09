package sns

// user gender
const (
	OAuthUserGenderFemale  int = iota // Female
	OAuthUserGenderMale               // male
	OAuthUserGenderUnknown            // unknown

	// oauth user platform
	OAuthUserPlatformWechat = "wechat" // wechat
	OAuthUserPlatformWeibo  = "weibo"  // weibo
	OAuthUserPlatformQQ     = "qq"     // qq
	OAuthUserPlatformGithub = "github" // github
)

// OAuthInterface interface
type OAuthInterface interface {
	Authorize(state, callback string) string
	AccessToken(code, state, callback string) (OAuthAccessToken, error)
	Userinfo(accessToken, openID string) (OAuthUserinfo, error)
}

// OAuthAccessToken access token
type OAuthAccessToken struct {
	AccessToken string `json:"access_token"`
	OpenID      string `json:"openid"`
	UnionID     string `json:"unionid"`
}

// OAuthUserinfo common struct
type OAuthUserinfo struct {
	Platform string `json:"platform"`
	Name     string `json:"name"`
	OpenID   string `json:"openid"`
	UnionID  string `json:"unionid"`
	Avatar   string `json:"avatar"`
	Gender   int    `json:"gender"`
	Detail   string `json:"-"`
}
