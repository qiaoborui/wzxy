package wzxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"wobuzaixiaoyuan/pkg/common"

	"github.com/pkg/errors"
)

type Area struct {
	ID        string `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Name      string `json:"name"`
	Radius    int    `json:"radius"`
	Shape     int    `json:"shape"`
}
type SignList struct {
	SignList []struct {
		SignId   string `json:"signId"`
		Id       string `json:"id"`
		AreaList []Area `json:"areaList"`
		SchoolId string `json:"schoolId"`
		UserArea string `json:"userArea"`
	}
}

type SignData struct {
	InArea     int     `json:"inArea"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	Street     string  `json:"street"`
	StreetCode string  `json:"streetcode"`
	AreaJson   string  `json:"areaJSON"`
	CityCode   string  `json:"citycode"`
	NationCode string  `json:"nationcode"`
	Adcode     string  `json:"adcode"`
	District   string  `json:"district"`
	Country    string  `json:"country"`
	Towncode   string  `json:"towncode"`
	Township   string  `json:"township"`
}

type ListResponse struct {
	Code int `json:"code"`
	Data []struct {
		AreaID     string `json:"areaId,omitempty"`
		AreaList   []Area `json:"areaList"`
		ID         string `json:"id"`
		SchoolID   string `json:"schoolId"`
		SignID     string `json:"signId"`
		SignStatus int    `json:"signStatus"`
		Type       int    `json:"type"`
		UserArea   string `json:"userArea"`
	} `json:"data"`
}

type SignResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s Session) GetSignList() (SignList, error) {
	resp, err := s.client.Get("https://gw.wozaixiaoyuan.com/sign/mobile/receive/getMySignLogs?page=1&size=10")
	if err != nil {
		return SignList{}, errors.Wrap(err, "发起请求失败")
	}
	defer resp.Body.Close()
	var data ListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return SignList{}, errors.Wrap(err, "解析响应失败")
	}
	result := SignList{}
	for _, v := range data.Data {
		if v.Type == 0 && v.SignStatus == 1 {
			result.SignList = append(result.SignList, struct {
				SignId   string `json:"signId"`
				Id       string `json:"id"`
				AreaList []Area `json:"areaList"`
				SchoolId string `json:"schoolId"`
				UserArea string `json:"userArea"`
			}{SignId: v.SignID, Id: v.ID, AreaList: v.AreaList, SchoolId: v.SchoolID, UserArea: v.UserArea})
		}

	}
	return result, nil
}

func (s Session) Sign() error {
	tasks, err := s.GetSignList()
	if err != nil {
		s.User.Result <- err.Error()
		return errors.Wrap(err, "获取签到任务失败")
	}
	if len(tasks.SignList) == 0 {
		s.User.Result <- fmt.Sprintf("[%s]没有签到任务", s.User.RealName)
		return nil
	}
	for _, v := range tasks.SignList {
		if len(v.AreaList) == 0 {
			continue
		}
		lng, err := common.ToFloat(v.AreaList[0].Longitude)
		if err != nil {
			return errors.Wrap(err, "经度转换失败")
		}
		lat, err := common.ToFloat(v.AreaList[0].Latitude)
		if err != nil {
			return errors.Wrap(err, "纬度转换失败")
		}
		jsonRequestBody := SignData{
			InArea:     1,
			Longitude:  lng,
			Latitude:   lat,
			Province:   "",
			City:       "",
			AreaJson:   AreaListToAreaJson(v.AreaList, v.UserArea),
			CityCode:   "",
			Adcode:     "",
			District:   "",
			Country:    "",
			Towncode:   "",
			Township:   "",
			StreetCode: "",
			Street:     "",
			NationCode: "",
		}
		requestData, _ := json.Marshal(jsonRequestBody)
		u := url.URL{
			Scheme: "https",
			Host:   "gw.wozaixiaoyuan.com",
			Path:   "sign/mobile/receive/doSignByArea",
		}
		query := u.Query()
		query.Set("id", v.Id)
		query.Set("signId", v.SignId)
		query.Set("schoolId", v.SchoolId)
		u.RawQuery = query.Encode()
		req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(requestData)))
		if err != nil {
			return errors.Wrap(err, "error creating new request")
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "error sending request")
		}
		defer resp.Body.Close()
		var data Response
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return errors.Wrap(err, "error decoding response")
		}

		if data.Code == 0 {
			s.User.Result <- fmt.Sprintf("[%s]签到成功", s.User.RealName)
		} else {
			s.User.Result <- fmt.Sprintf("[%s]签到失败,响应：%s", s.User.RealName, data.Message)
		}
	}

	return nil
}

func AreaListToAreaJson(areaList []Area, userArea string) string {
	for _, v := range areaList {
		if v.Name == userArea {
			Type := v.Shape
			latitude := v.Latitude
			longitude := v.Longitude
			radius := v.Radius
			id := v.ID
			name := v.Name
			result := fmt.Sprintf("{\"type\":%d,\"circle\":{\"latitude\":\"%s\",\"longitude\":\"%s\",\"radius\":%d},\"id\":\"%s\",\"name\":\"%s\"}", Type, latitude, longitude, radius, id, name)
			return result
		} else {
			fmt.Println(v.Name)
		}
	}
	return ""
}
