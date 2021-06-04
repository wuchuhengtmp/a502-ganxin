/**
 * @Desc    维修管理员角色集成测试
 * @Author  wuchuheng<wuchuheng@163.com>
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

// 仓库管理员测试上下文
var maintenanceAdminTestCtx = struct {
	Token string
	Username string
	Password string
}{
	Username: seeders.MaintenanceAdmin.Username,
	Password: seeders.MaintenanceAdmin.Password,
}

/**
 * 维修管理员登录测试
 */
func TestMaintenanceAdminRoleLogin(t *testing.T) {
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
		"phone": maintenanceAdminTestCtx.Username,
		"password": maintenanceAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleMaintenanceAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	maintenanceAdminTestCtx.Token = tokenInfo["accessToken"].(string)
}


/**
 * 维修管理员获取公司列表集成测试
 */
func TestMaintenanceAdminRoleGetAllCompany(t *testing.T)  {
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
	res, err := graphReqClient(q, v, roles.RoleMaintenanceAdmin)
	hasError(t, err )
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题") )
	}
}


/**
 * 维修管理员获取公司人员列表集成测试
 */
func TestMaintenanceAdminRoleGetCompanyUsers(t *testing.T)  {
	q := `
		query getCompanyUserQuery {
		  getCompanyUser{
			id
			role {
			  id
			  tag
			  name
			}
			phone
			wechat
			avatar{
			  id
			  url
			}
			isAble
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleMaintenanceAdmin)
	hasError(t, err)
}

/**
 * 维修管理员获取仓库列表集成测试
 */
func TestMaintenanceAdminRoleGetRepository(t *testing.T)  {
	q := `
		 query {
		  getRepositoryList {
			id
			isAble
			weight
			adminName
			adminWechat
			adminPhone
			address
			total
			isAble
			remark
			city
			pinYin
			name
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{} {},
	}
	_, err := graphReqClient(q, v, roles.RoleMaintenanceAdmin)
	hasError(t, err)
}
