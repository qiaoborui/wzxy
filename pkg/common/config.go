package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	setTime()
	// 设置配置文件名为 config
	viper.SetConfigName("config")
	// 设置配置文件格式为 JSON
	viper.SetConfigType("json")
	// 添加配置文件查找路径
	viper.AddConfigPath("./")

	// 读取配置文件
	_ = viper.ReadInConfig()

	// 从环境变量中读取配置信息，若环境变量中不存在该项则取配置文件中的值
	APPID = os.Getenv("APPID")
	if APPID == "" {
		APPID = viper.GetString("leancloud.appId")
	}

	APPKEY = os.Getenv("APPKEY")
	if APPKEY == "" {
		APPKEY = viper.GetString("leancloud.appKey")
	}

	APPURL = os.Getenv("APPURL")
	if APPURL == "" {
		APPURL = viper.GetString("leancloud.appUrl")
	}

	if APPID == "" || APPKEY == "" || APPURL == "" {
		fmt.Println("LeanCloud 配置信息不完整")
		os.Exit(1)
	} else {
		fmt.Println("LeanCloud 配置信息读取成功")
	}

}

var APPID string
var APPKEY string
var APPURL string
