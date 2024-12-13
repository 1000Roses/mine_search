package utils_call

import (
	"time"

	"github.com/go-resty/resty/v2"
)

func Request(url string, isGet bool, headers, param map[string]string, body map[string]interface{}, timeout int) (*resty.Response, error) {
	client := resty.New()
	if timeout != 0 {
		client.SetTimeout(time.Second * time.Duration(timeout))
	}
	if headers == nil {
		headers = map[string]string{}
	}
	if body == nil {
		body = map[string]interface{}{}
	}
	if param == nil {
		param = map[string]string{}
	}
	headers["Content-Type"] = "application/json"
	req := client.R().SetHeaders(headers).SetQueryParams(param).
		SetBody(body)
	if isGet {
		return req.Get(url)
	}
	return req.Post(url)
}
