package auth

import (
	"encoding/json"
	"errors"
	"github.com/Privasea/middleware/auth/utils"
)

type AdminAuth struct {
}

func (a AdminAuth) Check(userToken string, apiServerName string, apiMethod string, apiPath string) (bool, error) {
	flag, err := getApiDataForAdmin(apiServerName, apiPath, apiMethod, userToken)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func getApiDataForAdmin(server, path, method, userToken string) (bool, error) {

	client := utils.NewClient("http://exception-service:5000/api/v1/inner_use")

	req := getApiDataForAdminRequestData{
		Server:     server,
		Path:       path,
		Method:     method,
		WalletSign: userToken,
	}

	postResponse, err := client.Post("/admin/check", req)
	if err != nil {
		return false, err
	}

	// 解析 POST 请求的响应数据
	var response getApiDataForAdminResponse
	err = json.Unmarshal(postResponse, &response)
	if err != nil {
		return false, err
	}
	if response.Code != 0 {
		return false, errors.New(response.Msg)
	}
	return response.Data.Pass, nil

	return false, nil
}

type getApiDataForAdminRequestData struct {
	Server     string `json:"server"`
	Path       string `json:"path"`
	Method     string `json:"method"`
	WalletSign string `json:"walletSign"`
}

// Response 定义了响应的结构体
type getApiDataForAdminResponse struct {
	Code int    `json:"code"` // 状态码
	Data data   `json:"data"` // 数据
	Msg  string `json:"msg"`  // 错误消息或成功消息
}
type data struct {
	Pass bool `json:"pass"`
}
