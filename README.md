## middleware企业中间件聚集地
### auth包,旨在统一提供所有接口鉴权服务

- https://github.com/Privasea/middleware

#### 引入middleware

```
go get github.com/Privasea/middleware
import "github.com/Privasea/middleware"
```
#### 添加gin的全局中间件

```
// 新增auth.GinInterceptor中间件（放在tl.GinInterceptor后面）
r.Use(tl.GinInterceptor)
r.Use(auth.GinInterceptor())
```
### cors包,跨域
#### 添加gin的全局中间件
```
// 新增cors.CORSMiddleware中间件（放在tl.GinInterceptor前面）
r.Use(cors.CORSMiddleware())
r.Use(tl.GinInterceptor)
r.Use(auth.GinInterceptor())
```
### limit包,限流
#### 添加gin的全局中间件
```
//同ip整站访问限制
//rateLimit窗口时间访问次数
//rateLimitTTL窗口时间，比如 60*time.Second 60s
//rdb github.com/go-redis/redis/v8 redis实例
r.Use(limit.NewRateLimitMiddleware(rateLimit, rateLimitTTL, rdb))

//同ip单个接口访问限制，可用作防重复点击，用于接口层中间件
ex:
chainGroup.POST("/bind_poh", middleware.NewRateLimitForRepetition(rdb), Wrap(bootstrap.Chain.BindPoh))

//同ip单个接口访问限制，自定义参数，用于接口层中间件
ex:
chainGroup.POST("/bind_poh", limit.NewRateLimitForApi(rateLimit, rateLimitTTL, rdb), Wrap(bootstrap.Chain.BindPoh))

```