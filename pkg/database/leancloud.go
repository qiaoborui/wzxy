package database

import (
	"fmt"
	"sort"
	"wobuzaixiaoyuan/pkg/common"

	"github.com/leancloud/go-sdk/leancloud"
	"github.com/pkg/errors"
)

var Client *leancloud.Client

type Response struct {
	Results []User `json:"results"`
}

func GetUsers() ([]User, error) {
	var users []User
	err := Client.Class("InSchool").NewQuery().EqualTo("status", 1).Find(&users)
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
	err := Client.Class("logs").NewQuery().Order("-createdAt").Limit(100).Find(&res)
	if err != nil {
		return nil, errors.Wrap(err, "获取日志数据失败")
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i].CreatedAt.After(res[j].CreatedAt) {
			return true
		}
		return false
	})
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
