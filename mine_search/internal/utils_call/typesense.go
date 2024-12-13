package utils_call

import (
	"encoding/json"
	"mine/internal"
	"mine/internal/repositories"
	"mine/internal/settings"
	"mine/internal/utils"
	"strings"

	"go.uber.org/zap"
)

func TypeSenseSearchText(s *settings.AppSettings, repo *repositories.Repositories, document, query string) (map[string]interface{}, error) {
	url := internal.Domains.TypeSense + strings.ReplaceAll(internal.Endpoints.TypeSense.TextSearch, "{{document}}", document) + "?" + query
	headers := map[string]string{
		"Content-type":        "application/json",
		"X-TYPESENSE-API-KEY": s.Cfgs.TypeSenseKey,
	}
	input := map[string]interface{}{}
	internal.Log.Info("Call "+url, zap.Any("input", input), zap.Any("headers", headers))
	resp, err := utils.Request(url, true, headers, map[string]string{}, input, 10, false)
	if err != nil {
		internal.Log.Error("Call "+url, zap.Any("input", input), zap.Any("resp", resp), zap.Error(err))
		return nil, err
	}
	internal.Log.Info("Response", zap.Any("url", url), zap.Any("header", headers), zap.Any("input", input), zap.Any("response", resp.String()))
	if resp.StatusCode() != 200 {
		return nil, err
	}
	res := map[string]interface{}{}
	err = json.Unmarshal([]byte(resp.String()), &res)
	if err != nil {
		internal.Log.Error("Error Unmarshal", zap.Any("res", res), zap.Error(err))
		return nil, err
	}
	return res, nil
}
