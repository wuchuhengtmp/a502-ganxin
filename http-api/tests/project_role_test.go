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
	"github.com/stretchr/testify/assert"
	"http-api/app/models/codeinfo"
	"http-api/app/models/devices"
	"http-api/app/models/roles"
	"http-api/pkg/model"
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

/**
 * 项目管理员获取物流列表集成测试
 */
func TestProjectAdminRoleGetExpressList(t *testing.T) {
	q := `
		query getExpressListQuery {
		  getExpressList{
			id
			name
			isDefault
			remark
		  }
		}
	`
	v := map[string]interface{}{}
	res, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("failed:项目管理员获取物流列表集成测试")
	}
	me, _ := GetUserByToken(projectAdminTestCtx.Token)
	items := res["getExpressList"].([]interface{})
	for _,item := range items {
		express := item.(map[string]interface{})
		id := express["id"].(float64)
		record := codeinfo.CodeInfo{}
		model.DB.Model(&record).Where("id = ?", int64(id)).First(&record)
		assert.Equal(t, record.CompanyId, me.CompanyId)
	}
}

/**
 * 项目管理员获取价格集成测试
 */
func TestProjectAdminRoleGetPrice(t *testing.T) {
	q := `
		 query getPriceQuery {
		  getPrice
		}
	`
	v := map[string]interface{} {}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("failed:项目管理员获取集成测试")
	}
}

/**
 * 项目管理员登录设备集成测试
 */
func TestProjectAdminRoleLoginDevice(t *testing.T) {
	q := `
		mutation login($phone: String!, $password: String!, $mac: String!) {
		  login(phone: $phone, password: $password, mac: $mac) {
			accessToken
			expired
			role
			roleName
		  }
		}
	`
	v := map[string]interface{} {
		"phone": "13427969607",
		"password": "12345678",
		"mac": "123:1242:1242:12412",
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("failed:项目管理员登录设备集成测试")
	}
}
/**
 * 项目管理员获取设备列表集成测试
 */
func TestProjectAdminGetDeviceList(t *testing.T) {
	q := `
		query getDeviceListQuery {
		  getDeviceList{
			id
			userInfo{
			  id
			  name
			}
			mac
			isAble
		  }
		}
	`
	v := map[string]interface{} {}
	res, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("failed:项目管理员获取设备列表集成测试")
	}
	// 断言响应的数据就是用户的项目名下的
	me, _ := GetUserByToken(projectAdminTestCtx.Token)
	items := res["getDeviceList"].([]interface{})
	for _, item := range items {
		cItem := item.(map[string]interface{})
		id := cItem["id"].(float64)
		d := devices.Device{}
		err := d.GetDeviceSelfById(int64(id))
		assert.NoError(t, err)
		u, err := d.GetUser()
		assert.NoError(t, err)
		assert.Equal(t, u.CompanyId, me.CompanyId)
	}
}
/**
 * 项目管理员获取型钢列表集成测试
 */
func TestProjectAdminGetSteelList(t *testing.T) {
	q := `
	query getSteelListQuery ($input: PaginationInput!){
		getSteelList(input: $input) {
		list{
		  id
		  state
		  totalUsageRate
		  repository{
			id
			name
		  }
		  manufacturer{
			id
			name
		  }
		  materialManufacturer{
			id
			name
		  }
		  turnover
		  usageYearRate
		  totalUsageRate
		  producedDate
		  specifcation{
			id
			specification
		  }
		}
		total
	  }  
	}
	`
	v := map[string]interface{} {
		"input": map[string]interface {} {
			"page": 1,
			"pageSize": 10,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("项目管理员获取型钢列表集成测试")
	}
}
/**
 * 项目管理员设置密码集成测试
 */
func TestProjectAdminSetPassword(t *testing.T) {
	q := `
		mutation ($input: SetPasswordInput!) {
		  setPassword(input: $input)
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface {} {
			"password": "12345678",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("项目管理员设置密码集成测试")
	}
}
/**
 * 项目管理员获取我的信息集成测试
 */
func TestProjectAdminGetMyInfo(t *testing.T) {
	q := `
		query getMyInfoQuery{
		  getMyInfo{
			id
			name
			company{
			  id
			  name
			}
			phone
			role{
			  id
			  name
			  tag
			}
		  }
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface {} { },
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("项目管理员获取我的信息集成测试")
	}
}
