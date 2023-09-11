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
		Area       string `json:"area,omitempty"`
		AreaID     string `json:"areaId,omitempty"`
		City       string `json:"city,omitempty"`
		AreaList   []Area `json:"areaList"`
		Country    string `json:"country,omitempty"`
		Date       int64  `json:"date,omitempty"`
		District   string `json:"district,omitempty"`
		ID         string `json:"id"`
		Latitude   string `json:"latitude,omitempty"`
		Longitude  string `json:"longitude,omitempty"`
		Name       string `json:"name"`
		Province   string `json:"province,omitempty"`
		SchoolID   string `json:"schoolId"`
		SignID     string `json:"signId"`
		SignStatus int    `json:"signStatus"`
		Start      int64  `json:"start"`
		Street     string `json:"street,omitempty"`
		Township   string `json:"township,omitempty"`
		Type       int    `json:"type"`
		UserID     string `json:"userId"`
	} `json:"data"`
}

type SignResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (s *Session) GetSignList() SignList {
	if s.Err != nil {
		// 如果s.Err不为空，也可以在这里进行处理，例如记录日志
		return SignList{}
	}
	resp, err := s.client.Get("https://gw.wozaixiaoyuan.com/sign/mobile/receive/getMySignLogs?page=1&size=10")
	if err != nil {
		s.Err = errors.Wrap(err, "获取请求列表发起请求失败")
		// 在这里可以记录错误日志，但不会返回错误
		return SignList{}
	}
	defer resp.Body.Close()
	var data ListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		s.Err = errors.Wrap(err, "获取请求列表解析响应失败")
		// 在这里可以记录错误日志，但不会返回错误
		return SignList{}
	}
	result := SignList{}
	for _, v := range data.Data {
		if v.Type == 0 && v.SignStatus == 1 {
			result.SignList = append(result.SignList, struct {
				SignId   string `json:"signId"`
				Id       string `json:"id"`
				AreaList []Area `json:"areaList"`
			}{SignId: v.SignID, Id: v.ID, AreaList: v.AreaList})
		}
	}
	return result
}

func (s *Session) Sign(tasks SignList) {
	if s.Err != nil {
		return
	}
	if len(tasks.SignList) == 0 {
		s.User.Result <- fmt.Sprintf("[%s]没有签到任务", s.User.RealName)
		return
	}
	for _, v := range tasks.SignList {
		if len(v.AreaList) == 0 {
			continue
		}
		lng, err := common.ToFloat(v.AreaList[0].Longitude)
		if err != nil {
			s.Err = errors.Wrap(err, "经度转换失败")
			return
		}
		lat, err := common.ToFloat(v.AreaList[0].Latitude)
		if err != nil {
			s.Err = errors.Wrap(err, "纬度转换失败")
			return
		}
		jsonRequestBody := SignData{
			InArea:     1,
			Longitude:  lng,
			Latitude:   lat,
			Province:   "陕西省",
			City:       "西安市",
			AreaJson:   AreaListToAreaJson(v.AreaList),
			CityCode:   "156610100",
			Adcode:     "610118",
			District:   "鄠邑区",
			Country:    "中国",
			Towncode:   "610118006",
			Township:   "草堂街道",
			StreetCode: "94871608551973499",
			Street:     "关中环线",
			NationCode: "156",
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
		query.Set("schoolId", "19")
		u.RawQuery = query.Encode()
		req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer([]byte(requestData)))
		if err != nil {
			s.Err = errors.Wrap(err, "error creating new request")
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.client.Do(req)
		if err != nil {
			s.Err = errors.Wrap(err, "error sending request")
			return
		}
		defer resp.Body.Close()
		var data Response
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			s.Err = errors.Wrap(err,"error unmarshalling response body")
		}
		if data.Code == 0 {
			s.User.Result <- fmt.Sprintf("[%s]签到成功", s.User.RealName)
		} else {
			s.Err =  fmt.Errorf(fmt.Sprintf("签到失败,响应：%s", data.Message))
		}
	}
}

func AreaListToAreaJson(areaList []Area) string {
	for _, v := range areaList {
		if v.ID == "190002" {
			Type := v.Shape
			latitude := v.Latitude
			longitude := v.Longitude
			radius := v.Radius
			id := v.ID
			name := v.Name
			result := fmt.Sprintf("{\"type\":%d,\"circle\":{\"latitude\":\"%s\",\"longitude\":\"%s\",\"radius\":%d},\"id\":\"%s\",\"name\":\"%s\"}", Type, latitude, longitude, radius, id, name)
			return result
		} else {
			fmt.Println(v.ID)
		}
	}
	return ""
}
