# open login for Go(lang)

 - wechat
 - weibo
 - qq
 - github

```go
import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limtech/sns"
)

// const
const (
	ServerURL = "https://github.com"
	// SNS keys
	SNSWechatID     = "SNS_WECHAT_ID"
	SNSWechatSecret = "SNS_WECHAT_SECRET"
	SNSWeiboID      = "SNS_WEIBO_ID"
	SNSWeiboSecret  = "SNS_WEIBO_SECRET"
	SNSQQID         = "SNS_QQ_ID"
	SNSQQSecret     = "SNS_QQ_SECRET"
	SNSGithubID     = "SNS_GITHUB_ID"
	SNSGithubSecret = "SNS_GITHUB_SECRET"
)

// OAuth controller
// authorize /oauth/:sns/authorize
// callback  /oauth/:sns/callback
// revoke    /oauth/:sns/revoke
type OAuth struct {
	Platform string
	State    string
	Code     string
}

// NewOAuth new OAuth
func NewOAuth() *OAuth {
	return &OAuth{}
}

// CallbackURL return callback url
func (o *OAuth) CallbackURL() string {
	return fmt.Sprintf(`%s/oauth/%s/callback`, ServerURL, o.Platform)
}

// Authorize redirect to sns authorize url
func (o *OAuth) Authorize(c *gin.Context) {
	o.Platform = c.Param("platform")
	o.State = "authorize"
	if o.Platform == "" {
		c.Redirect(http.StatusFound, ServerURL)
		return
	}

	// redirect to new uri
	c.Redirect(
		http.StatusFound,
		o.Instance().(sns.OAuthInterface).Authorize(o.State, o.CallbackURL()),
	)
}

// Instance return sns.OAuthInterface
func (o *OAuth) Instance() sns.OAuthInterface {
	var oauth sns.OAuthInterface
	switch o.Platform {
	case sns.OAuthUserPlatformWechat:
		oauth = sns.NewOAuthWechat(SNSWechatID, SNSWechatSecret)
		break
	case sns.OAuthUserPlatformWeibo:
		oauth = sns.NewOAuthWeibo(SNSWeiboID, SNSWeiboSecret)
		break
	case sns.OAuthUserPlatformQQ:
		oauth = sns.NewOAuthQQ(SNSQQID, SNSQQSecret)
		break
	case sns.OAuthUserPlatformGithub:
		oauth = sns.NewOAuthGithub(SNSGithubID, SNSGithubSecret)
		break
	}
	return oauth
}
```