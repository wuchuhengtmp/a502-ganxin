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

/**
 * 项目管理员获取公司人员列表集成测试
 */
func TestProjectAdminRoleGetCompanyUsers(t *testing.T)  {
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
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	hasError(t, err)
}


/**
 * 项目管理员获取仓库列表集成测试
 */
func TestProjectAdminRoleGetRepository(t *testing.T)  {
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
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	hasError(t, err)
}

/**
 * 仓库项目员获取规格列表集成测试
 */
func TestProjectAdminRoleGetSpecification(t *testing.T) {
	q := `
		query getSpecificationQuery {
		  getSpecification {
			id
			type
			specification
			weight
			isDefault
			length
		  }
		}
	`
	v := map[string]interface{} {}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	hasError(t, err)
}

/**
 * 项目管理员获取材料商列表集成测试
 */
func TestProjectAdminRoleGetMaterialManufacturers(t *testing.T) {
	q := `
		query  getMaterialManufacturersQuery {
		  getMaterialManufacturers{
			id
			name
			isDefault
		  }
		}
	`
	v := map[string]interface{} {}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	hasError(t, err)
}

/**
 * 项目管理员获取制造商集成测试
 */
func TestProjectRoleGetManufacturer(t *testing.T) {
	q := `query getManufacturersQuery{
		getManufacturers{
		id
		isDefault
		name
		remark
	  }
	}`
	v := map[string]interface{} {}
	res, err := graphReqClient(q, v, roles.RoleProjectAdmin);
	if  err != nil {
		t.Fatal("failed:项目管理员获取制造商集成测试")
	}
	assertCompanyIdForGetManufacturers(t, res, projectAdminTestCtx.Token)
}
