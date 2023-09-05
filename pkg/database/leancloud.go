package database

import (
	"encoding/json"
	"fmt"
	"github.com/leancloud/go-sdk/leancloud"
	"github.com/pkg/errors"
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
	var users []User
	err := Client.Class("InSchool").NewQuery().Find(&users)
	if err != nil {
		fmt.Println(err)
		return nil, errors.Wrap(err, "获取用户数据失败")
	}
	return users, nil
}

func UploadLog(content string) error {
	log := Log{
		Content: content,
	}
	_, err := Client.Class("logs").Create(&log)
	return err
}

func GetLogs() ([]Log, error) {
	var res []Log
	err := Client.Class("logs").NewQuery().Find(&res)
	if err != nil {
		return nil, errors.Wrap(err, "获取日志数据失败")
	}
	return res, nil
}

func GetSpecficLog(id string) (Log, error) {
	var res Log
	err := Client.Class("logs").NewQuery().EqualTo("objectId", id).First(&res)
	if err != nil {
		return Log{}, errors.Wrap(err, "获取日志数据失败")
	}
	return res, nil
}

func init() {
	Client = leancloud.NewClient(&leancloud.ClientOptions{
		AppID:      common.APPID,
		AppKey:     common.APPKEY,
		ServerURL:  common.APPURL,
		Production: "0",
	})
}
