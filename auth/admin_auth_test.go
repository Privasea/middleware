package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用例结构体
type testCase struct {
	server    string
	path      string
	method    string
	userToken string
	address   string
	expected  bool
}

func TestGetApiDataForAdmin(t *testing.T) {
	// 测试数据
	testCases := []testCase{
		{"123", "222", "GET", "111", "222", false},
		{"testnet-service", "/api/v1/test", "POST", "0x89db0c3baca37b5b9d7f70c95fd4c511a953227e55b4ff7309be50703c94a29f030ed4cb20e8f1cc8f65660d0a338d965725a4880662e22f0b1da7460596e9061b", "0x9C33e9a4cC4C772ef5b36C664254c6dBBa10247D", true},
		{"testnet-service", "/api/v1/test", "POST", "0x89db0c3baca37b5b9d7f70c95fd4c511a953227e55b4ff7309be50703c94a29f030ed4cb20e8f1cc8f65660d0a338d965725a4880662e22f0b1da7460596e9061b", "0x9C33e9a4cC4C772ef5b36C664254c6dBBa1024", false},
	}

	// 循环遍历测试数据
	for _, tc := range testCases {
		t.Run(tc.userToken, func(t *testing.T) {

			result, _ := getApiDataForAdmin(tc.server, tc.path, tc.method, tc.userToken, tc.address)

			assert.Equal(t, tc.expected, result)
		})
	}
}
