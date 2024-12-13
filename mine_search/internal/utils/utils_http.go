package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func Request(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, timeout int, proxy bool) (*resty.Response, error) {
	now := GetTimeUTC7()
	defer func() { fmt.Printf("ExecTime url=%s, dt=%v \n", url, GetTimeUTC7().Sub(now).Milliseconds()) }()
	if headers == nil {
		headers = map[string]string{}
	}
	if param == nil {
		param = map[string]string{}
	}
	if body == nil {
		body = map[string]interface{}{}
	}
	_, ok := headers["Content-Type"]
	if !ok {
		headers["Content-Type"] = "application/json"
	}
	client := resty.New()
	client.SetTimeout(time.Second * time.Duration(timeout))
	if proxy {
		client.SetProxy("http://proxy:80")
	}
	req := client.R().SetHeaders(headers).SetQueryParams(param).
		SetBody(body)
	if isGet {
		return req.Get(url)
	}
	return req.Post(url)
}

func ResponseString(resp *resty.Response) interface{} {
	if resp == nil {
		return nil
	}
	result := map[string]interface{}{}
	err := json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return resp.String()
	}
	return result
}

func LogResponse(resp *resty.Response) []zap.Field {
	if resp == nil {
		return nil
	}
	result := []zap.Field{}
	result = append(result, zap.Any("http_status", resp.Status()), zap.Any("http_status", resp.Status()), zap.Any("resp", ResponseString(resp)), zap.Any("dt", fmt.Sprintf("%v", resp.Time().Seconds())))
	return result
}
