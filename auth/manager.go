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
func GinInterceptor(ctx *gin.Context, mockAuth Auth) {

	if utils.PathContainsKey(ctx.Request.URL.Path, "admin") {
		userToken := ctx.GetHeader("sign_data")
		apiServerName := ctx.GetHeader("server_name")
		address := ctx.GetHeader("address")
		apiMethod := ctx.Request.Method
		apiPath := ctx.Request.URL.Path
		algorithmAuth := newAuth("admin")
		if mockAuth != nil {
			algorithmAuth = mockAuth
		}
		flag, err := algorithmAuth.Check(userToken, apiServerName, apiMethod, apiPath, address)
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
