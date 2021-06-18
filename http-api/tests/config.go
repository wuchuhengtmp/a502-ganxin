/**
 * @Desc    测试配置
 * @Author  wuchuheng<root@wuchuheng.com>
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
	"http-api/app/models/codeinfo"
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/bootstrap"
	"http-api/config"
	pkgConfig "http-api/pkg/config"
	"http-api/pkg/jwt"
	"http-api/pkg/model"
	"testing"
)

var bashUrl string

// graphql 请求客户端
var client  *graphql.Client

func init() {
	config.Initialize()
	bootstrap.SetupDB()
	appPort := pkgConfig.GetString("APP_PORT")
	bashUrl = fmt.Sprintf("http://localhost:%s/query", appPort)
	client = graphql.NewClient(bashUrl)
}

/**
 * 断言错误并终止测试
 */
func hasError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

/**
* graphql 客户端
 */
func graphReqClient(query string, variables map[string]interface{}, role roles.GraphqlRole, p ...interface{}) (responseData map[string]interface{}, err error) {
	req := graphql.NewRequest(query)
	for key, variable := range variables {
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
	case roles.RoleCompanyAdmin:
		if len(companyAdminTestCtx.Token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", companyAdminTestCtx.Token))
		}
		break
	case roles.RoleRepositoryAdmin:
		if len(repositoryAdminTestCtx.Token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", repositoryAdminTestCtx.Token))
		}
		break
	case roles.RoleProjectAdmin:
		if len(p) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p[0]))
		} else if len(projectAdminTestCtx.Token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", projectAdminTestCtx.Token))
		}
		break
	case roles.RoleMaintenanceAdmin:
		if len(p) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p[0]))
		} else if len(maintenanceAdminTestCtx.Token) > 0 {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", maintenanceAdminTestCtx.Token))
		}
		break
	}

	ctx := context.Background()
	err = client.Run(ctx, req, &responseData)

	return responseData, err
}

/**
 * 通过token用获取户信息
 */
func GetUserByToken(token string) (*users.Users, error) {
	payload, _ := jwt.ParseByTokenStr(token)

	u := users.Users{}
	if err := model.DB.Model(&users.Users{}).Where("id = ?", payload.Uid).First(&u).Error; err != nil {
		return nil , err
	}

	return &u, nil
}

/**
 * 断言获取制造商的响应数据必须归属于对应的公司id
 */
func assertCompanyIdForGetManufacturers(t *testing.T, res map[string]interface{}, token string)  {
	me, _ := GetUserByToken(token)
	items := res["getManufacturers"].([]interface{})
	for _, item := range items {
		tmp := item.(map[string]interface{})
		id := tmp["id"].(float64)
		c := codeinfo.CodeInfo{}
		model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", int64(id)).First(&c)
		if c.CompanyId != me.CompanyId {
			t.Fatal("failed:当前数据项归属的公司与当前用户归属的公司不一致")
		}
	}
}
