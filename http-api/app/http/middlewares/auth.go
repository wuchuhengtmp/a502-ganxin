/**
 * 登录鉴权
 */
package middlewares

import (
	"http-api/pkg/jwt"
	"http-api/pkg/response"
	"net/http"
)

func Auth(next HttpHandlerFunc) HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		var errRes  =  response.Error{
			Errors: map[string][]string{
				"authorization": {"没有token, 鉴权失败"},
			},
			ErrorCode: response.ErrorCodes.InvalidToken,
		}
		if len(token)  <= 7 {
			errRes.ResponseByHttpWriter(w)
			return
		}
		tokenStr := token[7:]
		_, err := jwt.ParseByTokenStr(tokenStr)
		if err != nil {
			errRes.Errors = map[string][]string {
				"authorization": {"无效token, 鉴权失败"},
			}
			errRes.ResponseByHttpWriter(w)
			return
		}
		next(w, r)
	}
}