package database

import "github.com/leancloud/go-sdk/leancloud"

type User struct {
	leancloud.Object
	Username  string `json:"username"`
	RealName  string `json:"Zh_name"`
	Password  string `json:"passwd"`
	IsEnable  int    `json:"status"`
	Jwsession string `json:"cache"`
	Start     string `json:"start"`
	End       string `json:"end"`
}
