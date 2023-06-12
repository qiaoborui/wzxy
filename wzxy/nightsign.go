package wzxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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
		Area           string `json:"area,omitempty"`
		AreaID         string `json:"areaId,omitempty"`
		City           string `json:"city,omitempty"`
		Classes        string `json:"classes"`
		ClassesID      string `json:"classesId"`
		College        string `json:"college"`
		AreaList       []Area `json:"areaList"`
		Country        string `json:"country,omitempty"`
		CreateCollege  string `json:"createCollege"`
		CreateHead     string `json:"createHead"`
		CreateName     string `json:"createName"`
		Date           int64  `json:"date,omitempty"`
		Degree         string `json:"degree"`
		District       string `json:"district,omitempty"`
		End            int64  `json:"end"`
		Head           string `json:"head"`
		ID             string `json:"id"`
		IsRead         int    `json:"isRead"`
		Latitude       string `json:"latitude,omitempty"`
		LeaderSign     int    `json:"leaderSign"`
		Longitude      string `json:"longitude,omitempty"`
		Major          string `json:"major"`
		Mode           int    `json:"mode"`
		Name           string `json:"name"`
		Number         string `json:"number"`
		Phone          string `json:"phone"`
		Province       string `json:"province,omitempty"`
		QrCode         int    `json:"qrCode"`
		ReadDate       int64  `json:"readDate"`
		SchoolID       string `json:"schoolId"`
		SignContext    string `json:"signContext"`
		SignDay        string `json:"signDay,omitempty"`
		SignID         string `json:"signId"`
		SignMode       int    `json:"signMode"`
		SignStatus     int    `json:"signStatus"`
		SignTitle      string `json:"signTitle"`
		SignUserID     string `json:"signUserId,omitempty"`
		SignUserName   string `json:"signUserName,omitempty"`
		SignUserNumber string `json:"signUserNumber,omitempty"`
		SignUserType   string `json:"signUserType,omitempty"`
		Start          int64  `json:"start"`
		Street         string `json:"street,omitempty"`
		TargetID       string `json:"targetId"`
		TargetName     string `json:"targetName"`
		TargetType     int    `json:"targetType"`
		Teacher        string `json:"teacher"`
		TeacherID      string `json:"teacherId"`
		Township       string `json:"township,omitempty"`
		Type           int    `json:"type"`
		UserArea       string `json:"userArea"`
		UserID         string `json:"userId"`
		UserType       string `json:"userType"`
		Year           string `json:"year"`
	} `json:"data"`
}

func (s Session) GetSignList() (SignList, error) {
	resp, err := s.client.Get("https://gw.wozaixiaoyuan.com/sign/mobile/receive/getMySignLogs?page=1&size=10")
	if err != nil {
		fmt.Println(err)
		return SignList{}, err
	}
	defer resp.Body.Close()
	var data ListResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return SignList{}, err
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
	return result, nil
}

type SignResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func toFloat(s string) float64 {
	result, _ := strconv.ParseFloat(s, 64)
	return result
}

func (s Session) Sign() error {
	tasks, err := s.GetSignList()
	if err != nil {
		s.User.Result <- string(err.Error())
		return err
	}
	if len(tasks.SignList) == 0 {
		s.User.Result <- fmt.Sprintf("[%s]没有签到任务", s.User.RealName)
		return nil
	}
	for _, v := range tasks.SignList {
		if len(v.AreaList) == 0 {
			continue
		}
		jsonRequestBody := SignData{
			InArea:     1,
			Longitude:  toFloat(v.AreaList[0].Longitude),
			Latitude:   toFloat(v.AreaList[0].Latitude),
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

		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := s.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//resp.Body.Close()
		var data SignResponse
		err = json.NewDecoder(resp.Body).Decode(&data)
		if data.Code == 0 {
			s.User.Result <- fmt.Sprintf("[%s]签到成功", s.User.RealName)
		} else {
			s.User.Result <- fmt.Sprintf("[%s]签到失败", s.User.RealName)
		}
	}
	//mu.Unlock()

	return nil
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
