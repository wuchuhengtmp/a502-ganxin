/**
 * @Desc
 * @Author  仓库角色集成测试
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/3
 * @Listen  MIT
 */
package tests

import (
	"fmt"
	"http-api/app/models/roles"
	"http-api/seeders"
	"testing"
)

// 项目管理员测试上下文
var projectAdminTestCtx = struct {
	Token string
	Username string
	Password string
}{
	Username: seeders.ProjectAdmin.Username,
	Password: seeders.ProjectAdmin.Password,
}

/**
 * 项目管理员登录测试-后台
 */
func TestProjectAdminRoleLogin(t *testing.T) {
	query := `
		mutation login ($phone: String!, $password: String!) {
		  login(phone: $phone, password: $password) {
			accessToken
			expired
			role
			roleName
		  }
		}
	`
	variables :=  map[string]interface{} {
		"phone": projectAdminTestCtx.Username,
		"password": projectAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleProjectAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	projectAdminTestCtx.Token = tokenInfo["accessToken"].(string)
}

/**
 * 项目管理员获取公司列表集成测试
 */
func TestProjectAdminRoleGetAllCompany(t *testing.T)  {
	q := `query {
			  getAllCompany {
				id
				name
				symbol
				logoFile{
					id
				  url
				}
				backgroundFile {
				  id
				  url
				}
				isAble
				phone
				wechat
				startedAt
				endedAt
				adminName
				adminPhone
				adminWechat
				adminAvatar{
				  id
				  url
				}
				createdAt
			}
		}
	`
	v := map[string]interface{}{}
	res, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	hasError(t, err )
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题") )
	}
}