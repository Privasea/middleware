package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) Check(userToken, apiServerName, apiMethod, apiPath string) (bool, error) {
	args := m.Called(userToken, apiServerName, apiMethod, apiPath)
	return args.Bool(0), args.Error(1)
}

// 测试GinInterceptor函数
func TestGinInterceptor(t *testing.T) {
	// 创建一个Gin引擎实例
	r := gin.Default()

	// 创建一个mockAuth实例
	mockAuth := new(MockAuth)

	// 创建GinInterceptor中间件
	r.Use(func(ctx *gin.Context) {
		GinInterceptor(ctx, mockAuth) // 使用mock
	})

	// 测试用例 1: Token验证通过
	t.Run("valid token", func(t *testing.T) {
		// 构造请求
		req, _ := http.NewRequest("GET", "/admin/test", nil)
		req.Header.Set("sign_data", "valid_token")
		req.Header.Set("server_name", "test_server")

		// 模拟Check方法返回的值
		mockAuth.On("Check", "valid_token", "test_server", "GET", "/admin/test").Return(true, nil)

		// 记录响应
		w := performRequest(r, req)

		// 验证响应
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockAuth.AssertExpectations(t)
	})

	// 测试用例 2: Token验证失败
	t.Run("invalid token", func(t *testing.T) {
		// 构造请求
		req, _ := http.NewRequest("GET", "/admin/test", nil)
		req.Header.Set("sign_data", "invalid_token")
		req.Header.Set("server_name", "test_server")

		// 模拟Check方法返回的值
		mockAuth.On("Check", "invalid_token", "test_server", "GET", "/admin/test").Return(false, nil)

		// 记录响应
		w := performRequest(r, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "Auth Failed")
		mockAuth.AssertExpectations(t)
	})

	// 测试用例 3: 模拟内部错误
	t.Run("error in check", func(t *testing.T) {
		// 构造请求
		req, _ := http.NewRequest("GET", "/admin/test", nil)
		req.Header.Set("sign_data", "any_token")
		req.Header.Set("server_name", "test_server")

		// 模拟Check方法返回错误
		mockAuth.On("Check", "any_token", "test_server", "GET", "/admin/test").Return(false, fmt.Errorf("internal error"))

		// 记录响应
		w := performRequest(r, req)

		// 验证响应
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Body.String(), "interNet")
		mockAuth.AssertExpectations(t)
	})
}

// 执行请求
func performRequest(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}
