/**
 * @Desc    测试配置
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/1
 * @Listen  MIT
 */
package tests

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"http-api/app/models/roles"
	"testing"
)

const bashUrl string = "http://localhost:9501/query"
// graphql 请求客户端
var client = graphql.NewClient(bashUrl)
/**
 * 断言错误并终止测试
 */
func hasError(t *testing.T, err error) {
	if err != nil {
		t.Fatal( err.Error())
	}
}

/**
* graphql 客户端
 */
func graphReqClient(query string, variables map[string]interface{}, role roles.GraphqlRole) (responseData map[string]interface{}, err error) {
	req := graphql.NewRequest(query)
	for key,  variable := range variables {
		req.Var(key, variable)
	}
	req.Header.Set("Cache-Control", "no-cache")
	switch role {
	// 登记超级管理员角色token用于鉴权接口使用
	case roles.RoleAdmin:
		if len(superAdminTestCtx.SuperAdminToken) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", superAdminTestCtx.SuperAdminToken))
		}
		break
	}
	ctx := context.Background()
	err = client.Run(ctx, req, &responseData)


	return responseData, err
}
