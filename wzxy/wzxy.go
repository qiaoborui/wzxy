package wzxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type User struct {
	RealName string
	Username string
	Password string
	Result   chan string
}

type Session struct {
	client *http.Client
	User   *User
}

//var mu sync.Mutex

type loginResp struct {
	Code int
}

// NewWzxy返回一个用于打卡的客户端实例
func NewWzxy(user *User) *Session {
	jar, _ := cookiejar.New(nil)
	return &Session{
		client: &http.Client{Jar: jar},
		User:   user,
	}
}

func (s Session) Login() error {
	Url, _ := url.Parse("https://gw.wozaixiaoyuan.com/basicinfo/mobile/login/username")
	params := url.Values{}
	params.Set("username", s.User.Username)
	params.Set("password", s.User.Password)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, _ := http.NewRequest("POST", urlPath, strings.NewReader("{}"))
	//addHeaders(req)
	respRaw, _ := s.client.Do(req)
	defer respRaw.Body.Close()
	body, _ := io.ReadAll(respRaw.Body)
	var resp loginResp
	json.Unmarshal(body, &resp)
	if resp.Code != 0 {
		return fmt.Errorf("登录失败")
	}
	return nil
}
