package wzxy

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

func NewSession(user *User) *Session {
	jar, _ := cookiejar.New(nil)
	return &Session{
		client: &http.Client{Jar: jar},
		User:   user,
	}
}

func (s Session) Login() error {
	Url, err := url.Parse("https://gw.wozaixiaoyuan.com/basicinfo/mobile/login/username")
	if err != nil {
		return errors.Wrap(err, "error parsing URL")
	}
	params := url.Values{}
	params.Set("username", s.User.Username)
	params.Set("password", s.User.Password)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, err := http.NewRequest("POST", urlPath, strings.NewReader("{}"))
	if err != nil {
		return errors.Wrap(err, "error creating new request")
	}
	respRaw, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error sending request")
	}
	defer respRaw.Body.Close()
	body, err := io.ReadAll(respRaw.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}
	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return errors.Wrap(err, "error unmarshalling response body")
	}
	code, ok := resp["code"].(int)
	if !ok {
		return fmt.Errorf("error parsing code")
	}
	if code != 0 {
		return fmt.Errorf("login failed")
	}
	return nil
}
