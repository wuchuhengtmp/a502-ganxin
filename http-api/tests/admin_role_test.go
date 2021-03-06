/**
 * @Desc    超级管理员角色集成测试
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/1
 * @Listen  MIT
 */
package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"http-api/app/models/configs"
	"http-api/app/models/roles"
	"http-api/seeders"
	"testing"
	"time"
)

// 超级管理员测试上下文
var superAdminTestCtx = struct {
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
	variables := map[string]interface{}{
		"phone":    superAdminTestCtx.Username,
		"password": superAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleAdmin)
	assert.NoError(t, err)
	token := res["login"]
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
		"input": map[string]interface{}{
			"name":              "公司名1",
			"pinYin":            "GSM1",
			"symbol":            "这是公司宗旨",
			"logoFileId":        6,
			"backgroundFileId":  5,
			"isAble":            true,
			"phone":             "13427969140",
			"wechat":            "wc20030318",
			"startedAt":         "2021-12-18 18:00:00",
			"endedAt":           "2022-12-18 18:00:00",
			"adminName":         "公司管理员" + fmt.Sprintf("%d", time.Now().Unix())[4:],
			"adminPassword":     "12345678",
			"adminPhone":        "1342" + fmt.Sprintf("%d", time.Now().Unix())[3:], // mock phone number
			"adminAvatarFileId": 4,
			"adminWechat":       "wc20030318",
		},
	}
	res, err := graphReqClient(q, v, roles.RoleAdmin)
	assert.NoError(t, err)
	createCompany, _ := res["createCompany"].(map[string]interface{})
	id := createCompany["id"].(float64)
	superAdminTestCtx.DeleteCompanyId = int64(id)
}

/**
 * 超级管理员获取全部公司列表集成测试
 */
func TestSuperAdminRoleGetAllCompany(t *testing.T) {
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
	hasError(t, err)
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
		"input": map[string]interface{}{
			"id":                2,
			"name":              "2",
			"pinYin":            "3",
			"symbol":            "4",
			"logoFileId":        1,
			"backgroundFileId":  2,
			"isAble":            true,
			"phone":             seeders.CompanyAdmin.Username,
			"wechat":            "12345678",
			"startedAt":         "2021-12-31 00:00:00",
			"endedAt":           "2022-12-31 00:00:00",
			"adminName":         "username_change_test" + fmt.Sprintf("%d", time.Now().Unix())[6:],
			"adminPassword":     seeders.CompanyAdmin.Password,
			"adminAvatarFileId": 4,
			"adminPhone":        "13427969604",
			"adminWechat":       "wc20030318_change_wechat_" + fmt.Sprintf("%d", time.Now().Unix())[6:],
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
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

/**
 * 超级管理员获取概览集成测试
 */
func TestAdminRoleGetSummary(t *testing.T) {
	q := `
		query {
		  getSummary {
			#### 资产概况1 ###
			feeTotal #总价值(万元)
			weightTotal #型钢总量(吨)
			yearFeeTotal # 今年新增价值(万元)
			yearWeightTotal # 今年新增型钢(吨)
			
			#### 资产概况2 ###
			idleWeightTotal #闲置量(吨)
			leaseWeightTotal #租赁数量(吨)
			maintenanceTotal #维修数量
			scrapWeightTotal # 报废量(吨)

			### 最近最近 盘点 ###
			lossTotal #丢失数量
			maintenanceWeightTotal # 维修量(吨)
			projectTotal #项目总数
			weightTotal #总重量
			leaseTotal #总体租出  
		  }
		}
	`
	v := map[string]interface{}{
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

/**
 * 超级管理员获取仓库列表（用于仪表盘）集成测试
 */
func TestAdminRoleGetRepositoryListForDashboard(t *testing.T) {
	q := `
		query {
		   getRepositoryListForDashboard {
			id
			name
		  } 
		}
	`
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

/**
 * 超级管理员获取仓库列表（用于仪表盘）集成测试
 */
func TestAdminRoleGetSteelSummaryForDashboard(t *testing.T) {
	q := `
		query ($input: GetSteelSummaryForDashboardInput!){
		  getSteelSummaryForDashboard(input: $input){
			crappedPercent # 报废
			lostPercent # 丢失
			maintainingPercent # 维修中
			storedPercent # 在库
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"repositoryId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

///**
// * 超级管理员获取项目列表（用于仪表盘）集成测试
// */
//func TestAdminRoleGetProjectListForDashboard(t *testing.T) {
//	q := `
//		query {
//			  getProjectListForDashboard {
//				id
//				name
//			  }
//		}
//	`
//	_, err := graphReqClient(q, v, roles.RoleAdmin)
//	hasError(t, err)
//}

/**
 * 超级管理员获取项目型钢列表（用于仪表盘）集成测试
 */
func TestAdminRoleGetSteelForDashboard(t *testing.T) {
	q := `
		##
		# 更新时间和操作人字段不要 项目过滤也不要了
		##
		query ($input: GetSteelForDashboardInput){
		  getSteelForDashboard(input: $input) {
			list {
			  id
			  orderSpecification {
				order {
				  orderNo # 订单号
				  project {
					name # 项目名
				  }
				}
			  }
			  steel {
				repository {
				  name # 仓库名
				}
				specifcation {
				  specification # 规格
				}
				stateInfo {
				  state 
				  desc # 当前状态
				}
			  }
			  stateInfo { # 使用状态
				state 
				desc # 使用状态 
			  }
			  locationCode # 安装编码
			  enterWorkshopAt # 入场时间
			  outWorkshopAt # 出场时间
			  enterRepositoryAt # 归库时间
			   
			}
			total # 总量
		  }
		}
	`

	v = map[string]interface{}{
		"input": map[string]interface{}{
			"page":         1,
			"pageSize":     10,
			"repositoryId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

/**
 * 超级管理员忘记密码型钢列表集成测试
 */
func TestAdminRoleForgetPassword(t *testing.T) {
	q := `
		mutation ($input: GetCodeForForgetPasswordInput!){
		  createCode(input: $input) {
			key
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"phone": "13427969604",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}

/**
 * 超级管理员获取短信配置集成测试
 */
func TestAdminRoleGetSMSConfig(t *testing.T) {
	q := `
		query {
		  getSMSConfig {
			accessKey # 短信accesskey 
			accessSecretKey  # 短信accessSecreptkey 
			sign # 签名
			template # 模板
		  }
		}
	`
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	hasError(t, err)
}
/**
 * 超级管理员设置短信配置集成测试
 */
func TestAdminRoleSetSMSConfig(t *testing.T) {
	q := `
		mutation ($input: SetSMSConfigInput!) {
		  setSMSConfig(input: $input) {
			accessKey # 
			accessSecretKey 
			sign
			template
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{} {
			"accessKey":       configs.GetGlobalVal(configs.SMS_ACCESS_KEY),
			"accessSecretKey": configs.GetGlobalVal(configs.SMS_ACCESS_SECRET_KEY),
			"sign":            configs.GetGlobalVal(configs.SMS_SIGN),
			"template":        configs.GetGlobalVal(configs.SMS_TEMPLATECODE),
		},
	}
	_, err := graphReqClient(q, v, roles.RoleAdmin)
	assert.NoError(t, err)
}
