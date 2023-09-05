package common

import (
	"fmt"
	"time"
)

// CompareTime @Title CompareTime
// @Description 比较输入时间和当前时间的大小
// @Param inputTime string 输入时间
// @Success bool
func CompareTime(inputTime string) bool {
	//构造包含当前日期的字符串
	dateStr := time.Now().Format("2006-01-02") //使用Go语言规定的"2006-01-02"作为日期格式

	//将输入时间和日期信息组合成一个完整的时间字符串
	fullTimeStr := fmt.Sprintf("%s %s", dateStr, inputTime)

	//解析时间字符串为time类型，获取输入时间
	layout := "2006-01-02 15:04"
	t, err := time.ParseInLocation(layout, fullTimeStr, time.Local)
	if err != nil {
		fmt.Println(err)
		return false
	}
	//获取当前时间
	currentTime := time.Now()
	//比较当前时间和输入时间的差值
	if t.After(currentTime) {
		//fmt.Printf("%s is after current time %s", inputTime, currentTime.Format(layout))
		return true
	} else if t.Equal(currentTime) {
		//fmt.Printf("%s is equal to current time %s", inputTime, currentTime.Format(layout))
		return false
	} else {
		//fmt.Printf("%s is before current time %s", inputTime, currentTime.Format(layout))
		return false
	}
}
func setTime() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}
