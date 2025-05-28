package auth

import (
	"github.com/Privasea/middleware/auth/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func newAuth(algorithm string) Auth {
	switch algorithm {
	case "admin":
		return &AdminAuth{}
	default:
		return nil
	}
}
func GinInterceptor(ctx *gin.Context) {

	if utils.PathContainsKey(ctx.Request.URL.Path, "admin") {
		userToken := ctx.GetHeader("sign_data")
		apiServerName := ctx.GetHeader("server_name")
		apiMethod := ctx.Request.Method
		apiPath := ctx.Request.URL.Path
		flag, err := newAuth("admin").Check(userToken, apiServerName, apiMethod, apiPath)
		if err != nil {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"code":  300500,
				"msg":   "interNet" + err.Error(),
				"error": "",
			})
			ctx.Abort()
			return
		}
		if !flag {
			ctx.JSON(http.StatusOK, map[string]interface{}{
				"code":  300501,
				"msg":   "Auth Failed",
				"error": "",
			})
			ctx.Abort()
			return
		}

	}

	ctx.Next()
}
