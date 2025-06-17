package limit

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"testing"
	"time"
)

// 模拟API请求，测试并发访问
func TestRateApiLimit(t *testing.T) {
	// 创建并初始化 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.98.229.91:6379", // Redis 地址
		Password: "",                  // Redis 密码
		DB:       0,                   // Redis 使用的数据库编号
	})

	// 创建Gin引擎并加载中间件
	router := setupRouterForApi(rdb)

	// 启动Gin引擎
	go func() {
		if err := router.Run(":8080"); err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	// 模拟并发请求
	const concurrentRequests = 10
	done := make(chan bool, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(i int) {
			// 发送 HTTP GET 请求
			resp, err := http.Get("http://localhost:8080/test")
			if err != nil {
				t.Errorf("Request %d failed: %v", i, err)
			} else {
				var response Response
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					log.Fatalf("Error decoding response: %v", err)
				}

				if response.Code == 0 {
					fmt.Printf("Request %d: Success\n", i)
				} else {
					fmt.Printf("Request %d: Rate limit exceeded (Status: %d)\n", i, response.Code)

				}
			}
			done <- true
		}(i)
	}

	// 等待所有并发请求完成
	for i := 0; i < concurrentRequests; i++ {
		<-done
	}
}

// 测试服务器路由
func setupRouterForApi(rdb *redis.Client) *gin.Engine {
	r := gin.Default()
	r.Use(NewRateLimitForRepetition(rdb))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
		})
	})
	return r
}
