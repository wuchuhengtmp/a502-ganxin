/**
 * @Desc    公司管理员角色集成测试
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/1
 * @Listen  MIT
 */
package tests

import (
	"fmt"
	"http-api/app/models/roles"
	"http-api/seeders"
	"testing"
	"time"
)
// 超级管理员测试上下文
var companyAdminTestCtx = struct{
	// token 用于角色鉴权
	Token string
	// 用于删除的公司id
	DeleteCompanyId int64
	// 账号
	Username string
	// 密码
	Password string
}{
	Username: seeders.CompanyAdmin.Username,
	Password: seeders.CompanyAdmin.Password,
}

/**
 * 公司管理员登录集成测试-管理后台登录
 */
func TestCompanyAdminRoleLogin(t *testing.T) {
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
		"phone": companyAdminTestCtx.Username,
		"password": companyAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleCompanyAdmin)
	hasError(t, err)
	token := res["login"]
	//["accessToken"]
	tokenInfo := token.(map[string]interface{})
	companyAdminTestCtx.Token = tokenInfo["accessToken"].(string)
}

/**
 * 公司管理员获取全部公司列表集成测试
 */
func TestCompanyAdminRoleGetAllCompany(t *testing.T)  {
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
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err )
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题") )
	}
}

/**
 * 公司管理员修改公司集成测试
 */
func TestCompanyAdminRoleEditCompany(t *testing.T) {
	q := `
		mutation editMutation($input: EditCompanyInput!) {
		  editCompany(input: $input){
			id
			name
			pinYin
			symbol
			logoFile {id url}
			backgroundFile {id url}
			isAble
			phone
			wechat
			startedAt
			endedAt
			adminName
			createdAt
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{} {
			"id": 2,
			"name": "2",
			"pinYin": "3",
			"symbol": "4",
			"logoFileId": 1,
			"backgroundFileId": 2,
			"isAble": true,
			"phone": seeders.CompanyAdmin.Username,
			"wechat": "12345678",
			"startedAt": "2021-12-31 00:00:00",
			"endedAt": "2022-12-31 00:00:00",
			"adminName": "username_change_test_with_company_role",
			"adminPassword": seeders.CompanyAdmin.Password,
			"adminAvatarFileId": 4,
			"adminPhone": seeders.CompanyAdmin.Username,
			"adminWechat": "admin_wechat_change_test_with_company_role",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 添加公司人员集成测试
*/
func TestCompanyAdminRoleCreateCompanyUser(t *testing.T)  {
	q := `
		mutation createUserMutation($input: CreateCompanyUserInput!){
		  createCompanyUser(input: $input){
			id
			role {
			  id
			  name
			  tag
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
	v := map[string]interface{}{
		"input": map[string]interface{} {
			"name": "username _for_TesCreateCompanyUser",
			"phone": fmt.Sprintf("1342%d", time.Now().Unix())[0:11] ,
			"password": "12345678",
			"avatarId": 1,
			"role": "repositoryAdmin",
			"wechat": "wechat_for_testCreateCompanyUser",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}
