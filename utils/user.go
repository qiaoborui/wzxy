package utils

type User struct {
	Username  string `json:"username"`
	RealName  string `json:"zh_name"`
	Password  string `json:"passwd"`
	IsEnable  int    `json:"status"`
	Jwsession string `json:"cache"`
}
