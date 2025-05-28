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
r.Use(auth.GinInterceptor)
```