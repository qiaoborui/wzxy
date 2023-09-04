package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leancloud/go-sdk/leancloud"
	"io"
	"net/http"
	"net/url"
	"sort"
	"wobuzaixiaoyuan/pkg/common"
)

var Client *leancloud.Client

type Response struct {
	Results []User `json:"results"`
}

func FetchData(condition string) (Response, error) {
	client := &http.Client{}
	Url, _ := url.Parse(common.APPURL)
	params := url.Values{}
	params.Set("where", condition)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Add("X-LC-Id", common.APPID)
	req.Header.Add("X-LC-Key", common.APPKEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	sort.Slice(response.Results, func(i, j int) bool {
		if response.Results[i].IsEnable > response.Results[j].IsEnable {
			return true
		}
		return false
	})
	return response, nil
}

func GetUsers() ([]User, error) {
	user := User{}
	err := Client.Class("InSchool").NewQuery().EqualTo("status", 1).First(&user)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("获取用户数据失败")
	}
	return []User{}, nil
}

func Initial() {
	Client = leancloud.NewClient(&leancloud.ClientOptions{
		AppID:     common.APPID,
		AppKey:    common.APPKEY,
		ServerURL: common.APPURL,
	})
}
