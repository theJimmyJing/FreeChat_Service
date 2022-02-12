package middleware

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"fmt"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, token := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
		fmt.Println(token)
		if !ok {
			log.NewError("","GetUserIDFromToken false ", c.Request.Header.Get("token"))
			c.Abort()
			http.RespHttp200(c, constant.ErrParseToken, nil)
			return
		}
	}
}
