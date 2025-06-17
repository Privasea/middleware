package limit

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strings"
	"time"
)

// 接口防重复提交
func NewRateLimitForRepetition(redisClient *redis.Client) gin.HandlerFunc {
	rateLimit := 1
	rateLimitTTL := 3 * time.Second

	return newRateLimitMiddleware(rateLimit, rateLimitTTL, redisClient, "1")
}

// 单接口限流
func NewRateLimitForApi(rateLimit int, rateLimitTTL time.Duration, redisClient *redis.Client) gin.HandlerFunc {

	return newRateLimitMiddleware(rateLimit, rateLimitTTL, redisClient, "1")
}

// 整站限流
func NewRateLimit(rateLimit int, rateLimitTTL time.Duration, redisClient *redis.Client) gin.HandlerFunc {
	return newRateLimitMiddleware(rateLimit, rateLimitTTL, redisClient, "")
}

// NewRateLimitMiddleware 创建流量限制中间件
func newRateLimitMiddleware(rateLimit int, rateLimitTTL time.Duration, redisClient *redis.Client, redisKey string) gin.HandlerFunc {

	var ctx = context.Background()
	rdb := redisClient

	return func(c *gin.Context) {
		// 获取客户端 IP 或用户 ID 作为 key
		path := c.Request.URL.Path

		if !strings.Contains(path, "/job/") && !strings.Contains(path, "/inner_use/") {

			clientIP := c.ClientIP()
			key := redisKey
			// 生成 Redis 键名
			if redisKey == "" {
				key = fmt.Sprintf("rate_limit:%s", clientIP)
			} else {

				method := c.Request.Method
				key = fmt.Sprintf("rate_limit_api:%s:%s:%s", clientIP, path, method)
			}

			// 使用 Lua 脚本保证计数器操作是原子性的
			luaScript := `
			local current
			current = redis.call('GET', KEYS[1])
			if current then
				if tonumber(current) >= tonumber(ARGV[1]) then
					return -1
				end
				redis.call('INCR', KEYS[1])
			else
				redis.call('SET', KEYS[1], 1, 'EX', ARGV[2])
                current = 1
			end
			return current
		`

			// 执行 Lua 脚本，返回值为当前请求次数
			result, err := rdb.Eval(ctx, luaScript, []string{key}, rateLimit, int(rateLimitTTL.Seconds())).Result()
			if err != nil {
				c.JSON(http.StatusOK, map[string]interface{}{
					"code":  200500,
					"msg":   "bad request",
					"error": "Error checking rate limit",
				})
				c.Abort()

				return
			}

			// 如果返回值为 -1，则表示超过请求限制
			if result == int64(-1) {
				c.JSON(http.StatusOK, map[string]interface{}{
					"code":  200429,
					"msg":   "Too many requests, please try again later",
					"error": "Too many requests, please try again later",
				})
				c.Abort()

				return
			}

			// 继续执行请求
			c.Next()
		}
	}
}
