package main

import (
	"fmt"
	"log"
	"os"
	"wobuzaixiaoyuan/pkg/common"
	"wobuzaixiaoyuan/pkg/database"
	"wobuzaixiaoyuan/pkg/wzxy"

	"github.com/jinzhu/copier"
)

type Instance struct {
	EventMap map[string]*struct {
		Count  int // <= 2
		Status int // 0: success, 1: failed
	}
	Users []database.User
}

func NewInstance() (*Instance, error) {
	users, err := database.GetUsers()
	if err != nil {
		return nil, err
	}
	return &Instance{
		Users: users,
		EventMap: make(map[string]*struct {
			Count  int // <= 2
			Status int // 0: success, 1: failed
		}),
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
	for _, user := range i.Users {
		// 如果当前时间在用户设置的时间范围内，并且用户今天还没有打过卡
		if !common.CompareTime(user.Start) && common.CompareTime(user.End) && i.EventMap[user.RealName].Count < 2 && i.EventMap[user.RealName].Status == 0 {
			wzxyUser := &wzxy.User{}
			copier.Copy(&wzxyUser, &user)
			checkinUsers = append(checkinUsers, wzxyUser)
			i.EventMap[user.RealName].Count += 1
		}
	}
	if len(checkinUsers) == 0 {
		log.SetOutput(os.Stdout)
		log.Printf("没有需要打卡的用户")
	}
	if len(checkinUsers) != 0 {
		res := wzxy.DoWork(checkinUsers)
		for _, result := range res {
			if result.Status == 0 {
				i.EventMap[result.RealName].Status = 0
			} else {
				i.EventMap[result.RealName].Status = 1
			}
		}
	}
}

func (i *Instance) ResetEventMap() {
	log.Printf("重置打卡次数")
	for _, user := range i.Users {
		i.EventMap[user.RealName].Count = 0
	}
}
