package main

import (
	"fmt"
	"log"
	"os"
	"wobuzaixiaoyuan/pkg/common"
	"wobuzaixiaoyuan/pkg/database"
	"wobuzaixiaoyuan/pkg/wzxy"
)

type Instance struct {
	EventMap map[string]int
	Users    []database.User
}

func NewInstance() (*Instance, error) {
	users, err := database.GetUsers()
	if err != nil {
		return nil, err
	}
	return &Instance{
		Users:    users,
		EventMap: make(map[string]int),
	}, nil
}

func (i *Instance) UpdateData() {
	log.Printf("更新用户数据")
	dataTmp, err := database.GetUsers()
	if err != nil {
		fmt.Println(err)
	} else {
		i.Users = dataTmp
	}
}

func (i *Instance) CheckInTask() {
	var checkinUsers []*wzxy.User
	results := make(chan string, 10)
	for _, user := range i.Users {
		// 如果当前时间在用户设置的时间范围内，并且用户今天还没有打过卡
		if !common.CompareTime(user.Start) && common.CompareTime(user.End) && i.EventMap[user.RealName] < 2 {
			checkinUsers = append(checkinUsers, &wzxy.User{
				RealName: user.RealName,
				Username: user.Username,
				Password: user.Password,
				Result:   results,
			})
			i.EventMap[user.RealName]++
		}
	}
	if len(checkinUsers) == 0 {
		log.SetOutput(os.Stdout)
		log.Printf("没有需要打卡的用户")
	}
	if len(checkinUsers) != 0 {
		wzxy.DoWork(checkinUsers)
	}
}

func (i *Instance) ResetEventMap() {
	log.Printf("重置打卡次数")
	for _, user := range i.Users {
		i.EventMap[user.RealName] = 0
	}
}
