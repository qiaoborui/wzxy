package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type Response struct {
	Results []User `json:"results"`
}

func FetchData(condition string) (Response, error) {
	client := &http.Client{}
	Url, _ := url.Parse(APPURL)
	params := url.Values{}
	params.Set("where", condition)
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	req, _ := http.NewRequest("GET", urlPath, nil)
	req.Header.Add("X-LC-Id", APPID)
	req.Header.Add("X-LC-Key", APPKEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Response{}, err
	}
	sort.Slice(response.Results, func(i, j int) bool {
		if response.Results[i].IsEnable > response.Results[j].IsEnable {
			return true
		}
		return false
	})
	return response, nil
}

func DisableUsers(id string) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("PUT", APPURL+id, strings.NewReader("{\"status\":0}"))
	req.Header.Add("X-LC-Id", APPID)
	req.Header.Add("X-LC-Key", APPKEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)
	if err == nil {
		return true
	}
	return false
}
