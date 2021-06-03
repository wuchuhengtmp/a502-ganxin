/**
 * @Desc    超级管理员角色集成测试
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
var superAdminTestCtx = struct{
	// token 用于角色鉴权
	SuperAdminToken string
	// 用于删除的公司id
	DeleteCompanyId int64
	// 账号
	Username string
	// 密码
	Password string
}{
	Username: seeders.SuperAdmin.Username,
	Password: seeders.SuperAdmin.Password,
}

/**
 * 超级管理员登录集成测试
 */
func TestSuperAdminRoleLogin(t *testing.T) {
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
		"phone": superAdminTestCtx.Username,
		"password": superAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleAdmin)
	hasError(t, err)
	token := res["login"]
	//["accessToken"]
	tokenInfo := token.(map[string]interface{})
	superAdminTestCtx.SuperAdminToken = tokenInfo["accessToken"].(string)
}

/**
 * 超级管理员创建公司集成测试
 */
func TestSuperAdminRoleCreateCompany(t *testing.T) {
	q := `
		mutation crateCompanyMutation($input: CreateCompanyInput!){
		  createCompany(input: $input){
			id
			name
			pinYin
			symbol
			logoFile {
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
	v := map[string]interface{}{
		"input": map[string]interface{} {
			"name": "公司名1",
			"pinYin": "GSM1",
			"symbol": "这是公司宗旨",
			"logoFileId": 6,
			"backgroundFileId": 5,
			"isAble": true,
			"phone": "13427969140",
			"wechat": "wc20030318",
			"startedAt": "2021-12-18 18:00:00",
			"endedAt": "2022-12-18 18:00:00",
			"adminName": "公司管理员" + fmt.Sprintf("%d", time.Now().Unix())[4:],
			"adminPassword": "12345678",
			"adminPhone": "1342" + fmt.Sprintf("%d", time.Now().Unix())[3:], // mock phone number
			"adminAvatarFileId": 4,
			"adminWechat": "wc20030318",
		},
	}
	res, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
	createCompany, _ :=  res["createCompany"].(map[string]interface{})
	id := createCompany["id"].(float64)
	superAdminTestCtx.DeleteCompanyId = int64(id)
}

/**
 * 超级管理员获取全部公司列表集成测试
 */
func TestSuperAdminRoleGetAllCompany(t *testing.T)  {
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
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err )
}

/**
 * 超级管理员修改公司集成测试
 */
func TestAdminRoleEditCompany(t *testing.T) {
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
			"adminName": "username_change_test" + fmt.Sprintf("%d", time.Now().Unix())[6:],
			"adminPassword": seeders.CompanyAdmin.Password,
			"adminAvatarFileId": 4,
			"adminPhone": "13427969604",
			"adminWechat": "wc20030318_change_wechat_" + fmt.Sprintf("%d", time.Now().Unix())[6:],
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

/**
 * 超级管理员删除公司集成测试
 */
func TestAdminRoleDeleteCompany(t *testing.T) {
	q := `
		mutation DeleteCompanyMutation($deleteCompanyId: Int!) {
		  deleteCompany(id: $deleteCompanyId)
		}
	`
	v := map[string]interface{}{
		"deleteCompanyId": superAdminTestCtx.DeleteCompanyId,
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}