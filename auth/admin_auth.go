package auth

import (
	"encoding/json"
	"errors"
	"github.com/Privasea/middleware/auth/utils"
)

type AdminAuth struct {
}

func (a AdminAuth) Check(userToken string, apiServerName string, apiMethod string, apiPath string, address string) (bool, error) {
	flag, err := getApiDataForAdmin(apiServerName, apiPath, apiMethod, userToken, address)
	if err != nil {
		return false, err
	}
	return flag, nil
}

func getApiDataForAdmin(server, path, method, userToken, address string) (bool, error) {

	//client := utils.NewClient("http://exception-service:5000/api/v1/inner_use")
	client := utils.NewClient("https://api-dev.privasea.ai/cloud-plat/api/v1/inner_use")

	req := getApiDataForAdminRequestData{
		ServiceName: server,
		Path:        path,
		Method:      method,
		SignData:    userToken,
		Address:     address,
	}

	postResponse, err := client.Post("/user/user_permission", req)
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
	return response.Data.Permission, nil
}

type getApiDataForAdminRequestData struct {
	ServiceName string `json:"service_name"`
	Path        string `json:"path"`
	Method      string `json:"method"`
	SignData    string `json:"sign_data"`
	Address     string `json:"address"`
}

// Response 定义了响应的结构体
type getApiDataForAdminResponse struct {
	Code int    `json:"code"` // 状态码
	Data data   `json:"data"` // 数据
	Msg  string `json:"msg"`  // 错误消息或成功消息
}
type data struct {
	Permission bool `json:"permission"`
}
