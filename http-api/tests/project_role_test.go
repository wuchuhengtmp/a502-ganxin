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
	"http-api/app/models/msg"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/roles"
	"http-api/app/models/specificationinfo"
	"http-api/pkg/model"
	"http-api/seeders"
	"testing"
	"time"
)

// 项目管理员测试上下文
var projectAdminTestCtx = struct {
	Token       string
	Username    string
	Password    string
	DeviceToken string
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
	variables := map[string]interface{}{
		"phone":    projectAdminTestCtx.Username,
		"password": projectAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleProjectAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	projectAdminTestCtx.Token = tokenInfo["accessToken"].(string)
}

/**
 * 项目管理员登录测试-手机机
 */
func TestProjectAdminRoleDeviceLogin(t *testing.T) {
	query := `
		mutation ($phone: String!, $password: String!, $mac: String!){
		  login (phone: $phone, password: $password, mac: $mac) {
			accessToken
		  }
		}
	`
	variables := map[string]interface{}{
		"phone":    projectAdminTestCtx.Username,
		"password": projectAdminTestCtx.Password,
		"mac":      "123:1242:1242:12412",
	}
	res, err := graphReqClient(query, variables, roles.RoleProjectAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	projectAdminTestCtx.DeviceToken = tokenInfo["accessToken"].(string)
}

/**
 * 项目管理员获取公司列表集成测试
 */
func TestProjectAdminRoleGetAllCompany(t *testing.T) {
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
	hasError(t, err)
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题"))
	}
}

/**
 * 项目管理员获取公司人员列表集成测试
 */
func TestProjectAdminRoleGetCompanyUsers(t *testing.T) {
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
func TestProjectAdminRoleGetRepository(t *testing.T) {
	q := `
		 query {
		  getRepositoryList {
			id
			isAble
			weight
			leaders {
				id
				name
				wechat
				phone
			}
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
		"input": map[string]interface{}{},
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
	v := map[string]interface{}{}
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
	v := map[string]interface{}{}
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
	v := map[string]interface{}{}
	res, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
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
	for _, item := range items {
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
	v := map[string]interface{}{}
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
	v := map[string]interface{}{
		"phone":    "13427969607",
		"password": "12345678",
		"mac":      "123:1242:1242:12412",
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
	v := map[string]interface{}{}
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
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"page":     1,
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
	v := map[string]interface{}{
		"input": map[string]interface{}{
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
	v := map[string]interface{}{
		"input": map[string]interface{}{},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	if err != nil {
		t.Fatal("项目管理员获取我的信息集成测试")
	}
}

/**
 * 项目管理员获取项目列表集成测试
 */
func TestProjectAdminGetProjectList(t *testing.T) {
	q := `
		query {
		  getProjectLis {
			id
			name
			address
			remark
			startedAt
			leaderList {
			  id
			  name
			}
			city
			endedAt
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取仓库概览测试
 */
func TestProjectAdminGetRepositoryOverview(t *testing.T) {
	me, _ := GetUserByToken(projectAdminTestCtx.Token)
	r := repositories.Repositories{}
	err := model.DB.Model(&repositories.Repositories{}).Where("company_id = ?", me.CompanyId).First(&r).Error
	assert.NoError(t, err)
	s := specificationinfo.SpecificationInfo{}
	err = model.DB.Model(&specificationinfo.SpecificationInfo{}).Where("company_id = ?", me.CompanyId).First(&s).Error
	assert.NoError(t, err)
	q := `
		query ($input: GetRepositoryOverviewInput!){
		  getRepositoryOverview(input: $input){
			total
			 weight
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":              r.ID,
			"specificationId": s.ID,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员创建需求单集成测试
 */
func TestProjectAdminCreateOrder(t *testing.T) {
	q := `
		mutation ($input: CreateOrderInput!) {
		  createOrder(input: $input) {
			id
			state
		  }
		}
	`
	me, _ := GetUserByToken(projectAdminTestCtx.Token)
	// 预计归还时间
	expectedReturnAt := time.Unix(time.Now().Unix()+60*60*24*30, 0).Format(time.RFC3339)
	// 项目id
	var ps []*projects.Projects
	err := model.DB.Model(&projects.Projects{}).
		Select(fmt.Sprintf("%s.*", projects.Projects{}.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", project_leader.ProjectLeader{}.TableName(), project_leader.ProjectLeader{}.TableName(), projects.Projects{}.TableName())).
		Where(fmt.Sprintf("%s.company_id = %d", projects.Projects{}.TableName(), me.CompanyId)).
		Where(fmt.Sprintf("%s.uid = %d", project_leader.ProjectLeader{}.TableName(), me.Id)).
		Find(&ps).
		Error
	assert.NoError(t, err)
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"expectedReturnAt": expectedReturnAt,
			"partList":         "这是配件清单_for_ProjectRoleCreateTest",
			"projectId":        ps[0].ID,
			"repositoryId":     1,
			"remark":           "这是备注",
			"steelList": []interface{}{
				map[string]interface{}{
					"total":           1,
					"specificationId": 1,
				},
			},
		},
	}
	_, err = graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取订单列表集成测试
 */
func TestProjectAdminGetOrderList(t *testing.T) {
	q := `
 query ($input: GetOrderListInput!){
		  getOrderList(input: $input) {
		   id
			state
			orderNo
			project {
			  id
			  name
			}
			repository{
			  id
			  name
			}
			state
			expectedReturnAt
			partList
			createdAt
			createUser {
			  id
			  name
			}
			confirmedUser {
			  id
			  name
			}
			confirmedAt
			expressList {
        id
        createdAt
        sender {
          id
          name
        }
        receiver {
          id
          name
        }
        direction
        expressCompany {
          id
          name
        }
        expressNo
        receiveAt
         
      }
			
			total
			weight		
			orderNo
			remark
		  }
		}

	`
	v := map[string]interface{}{
		"input": map[string]interface{}{},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取订单列表集成测试-手持机
 */
func TestProjectAdminDeviceGetOrderList(t *testing.T) {
	q := `
		 query ($input: GetOrderListInput!){
		  getOrderList(input: $input) {
		   id
			state
			orderNo
			project {
			  id
			  name
			}
			repository{
			  id
			  name
			}
			state
			expectedReturnAt
			partList
			createdAt
			createUser {
			  id
			  name
			}
			confirmedUser {
			  id
			  name
			}
			confirmedAt
			expressList {
        id
        createdAt
        sender {
          id
          name
        }
        receiver {
          id
          name
        }
        direction
        expressCompany {
          id
          name
        }
        expressNo
        receiveAt
         
      }
			
			total
			weight		
			orderNo
			remark
		  }
		}

	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"queryType": "toBeConfirm",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取订单详情集成测试-手持机
 */
func TestProjectAdminDeviceGetOrderDetail(t *testing.T) {
	me, err := GetUserByToken(projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
	o := orders.Order{}
	err = model.DB.Model(&orders.Order{}).Where("company_id = ?", me.CompanyId).First(&o).Error
	assert.NoError(t, err)
	q := `
		query ($input: getOrderDetailInput!){
		  getOrderDetail(input: $input) {
			id
			project {
			  id
			  name
		   }
			repository {
			  id
			  name
			}
			remark
			partList
			orderSpecificationList{
			  id
			  specification
			  total
			  weight
			}
			state
			expectedReturnAt
		   createdAt
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id": o.Id,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员型钢快速查询集成测试-手持机
 */
func TestProjectAdminGetSteelDetail(t *testing.T) {
	q := `query ($input: GetOneSteelDetailInput!){
		  getOneSteelDetail(input: $input) {
			id
			identifier
			specifcation{
			  specification
			}
			state
			manufacturer {
			  id
			  name
			}
			manufacturer {
			  id
			  name
			}
			producedDate
			createdAt
			createUser{
			  id
			  name
			}
			steelInProject{
			  id
			  outRepositoryAt
			  order{
				id
			  }
			  locationCode
			   useDays
			  order {
				project{
				  name
				}
			  }
			  enterWorkshopAt
			  outWorkshopAt
			  installationAt
			  outRepositoryAt
			  enterRepositoryAt
			}
			steelInMaintenance {
			  id
			  outedAt
			  maintenance{
				name
			  }
			  useDays
			  enteredAt
			  outedAt
			  enterRepositoryAt
			}
		  }
		} 
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifier": "8",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取入场订单详情集成测试
 */
func TestProjectAdminRoleGetSend2WorkshopListDetail(t *testing.T) {
	q := `
		query ($input: GetProjectOrder2WorkshopDetailInput!){
		  getSend2WorkshopOrderListDetail(input: $input) {
			list{
			  id
			}
			total
			totalWeight
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"orderId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.Error(t, err)
}

/**
 * 项目管理员型钢入场集成测试
 */
func testProjectAdminRoleSetProjectEnterWorkshop(t *testing.T) {
	orderId := 1
	q := `
			mutation ($input: SetSteelIntoWorkshopInput!) {
			  setSteelEnterWorkshop(input: $input){
				id
			  }
			}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"orderId": orderId,
			"identifierList": []interface{}{
				"8",
			},
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取项目规格列表集成测试--手机
 */
func testProjectAdminRoleGetProjectSpecificationDetail(t *testing.T) {
	q := `
		query ($input: GetProjectSpecificationDetailInput!) {
		  getProjectSpecificationDetail(input: $input) {
			list {
			  id
			  specificationInfo{
				specification
				id
				weight
			  }
			  # 应该接收(总量) 注: 重量 自己用这个数量乘规格信息中的重量就可以了
			  total
			  # 已接收
			  workshopReceiveTotal
			  # 已归库
			  storeTotal
			}
			total
			weight
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取项目详情列表集成测试--手机
 */
func testProjectAdminRoleGetProjectSteelDetail(t *testing.T) {
	q := `
		query ($input: GetProjectSteelDetailInput!){
		  getProjectSteelDetail(input: $input) {
			list{
			  id
			  # 订单规格列表
			  orderSpecification{
				# 订单
				order {
				  # 订单号
				  id
				}
			  }
			  # 位置编码
			  locationCode
			  # 型钢
			  steel{
				# 型钢编码
				code
				# 规格尺寸
				specifcation{
				  specification
				}
				# 周转次数
				turnover
			  }
			  # 出库时间
			  outWorkshopAt
			  #入场时间
			  enterWorkshopAt
			  # 出场时间
			  outWorkshopAt
			}
			total
			weight
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId":       1,
			"state":           200,
			"specificationId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取项目型钢状态列表集成测试--手机
 */
func testProjectAdminRoleGetProjectSteelStateList(t *testing.T) {
	q := `
		query {
		  # 获取项目型钢状态列表
		  getProjectSteelStateList{
			# 状态
			state
			# 说明
			desc
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取项目型钢状态列表集成测试--手机
 */
func testProjectAdminRoleGetMaxLocationCode(t *testing.T) {
	q := `
		query ($input: GetMaxLocationCodeInput!){
		  getMaxLocationCode(input: $input)
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifier": "8",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员安装码是否可用集成测试--手持机
 */
func testProjectAdminRoleIsAccessLocationCode(t *testing.T) {
	q := `
		query ($input: IsAccessLocationCodeInput!){
		  isAccessLocationCode(input: $input) 
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifier":   "8",
			"locationCode": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员安装型钢集成测试-手持机
 */
func testProjectAdminRoleInstallSteel(t *testing.T) {
	q := `
		mutation ($input: InstallLocationInput!) {
		  installSteel(input: $input)
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifier":   "8",
			"locationCode": 2,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 *项目管理员获取消息列表集成测试-手持机
 */
func TestProjectAdminRoleGetMsgList(t *testing.T) {
	q := `
		query {
		  getMsgList{
			id
			isRead
			content
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 *项目管理员待修改型钢详细信息集成测试-手持机
 */
func testProjectAdminRoleGetProjectSteel2BeChangeDetail(t *testing.T) {
	q := `
		query ($input: ProjectSteel2BeChangeInput!) {
		  getProjectSteel2BeChangeDetail(input: $input) {
			list {
			   id
				# 需求单号
				orderSpecification {
				  order{
					orderNo
				  }
				}
				# 当前状态
				stateInfo{
				  state
				  desc
				}
				# 规格尺寸
				orderSpecification {
				  specificationInfo{
					specification
				  }
				}
				# 周转次数
				steel{
				  turnover
				}
			}
			
			#数量
			total
			# 重量
			weightTotal
			}
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifierList": []string{
				"8",
			},
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员 项目中的型钢查询 集成测试-手持机
 */
func testProjectAdminRoleGetProjectSteel2BeChange(t *testing.T) {
	q := `
		query ($input:  GetProjectSteel2BeChangeInput!){
		  getProjectSteel2BeChange(input: $input) {
			id
			# 规格
			orderSpecification{
			  specificationInfo{
				specification
				id 
				# 重量
				weight
			  }
			}
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifier": "8",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员 修改项目型钢状态 集成测试-手持机
 */
func testProjectAdminRoleSetProjectSteelState(t *testing.T) {
	q := `
		mutation ($input: SetProjectSteelInput!){
		  setProjectSteelState(input: $input)
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifierList": []string{
				"8",
			},
			"state": 200,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 基础上管理员获取可出场的项目列表集成测试--手持机
 */
func testProjectAdminRoleGetOutOfWorkshopProjectList(t *testing.T) {
	q := `
		query {
		  getOutOfWorkshopProjectList {
			id
			name
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取可出场的项目列表集成测试--手持机
 */
func testProjectAdminRoleSetProjectSteelOutOfWorkshop(t *testing.T) {
	q := `
		mutation ($input: SetProjectSteelOutOfWorkshopInput!){
		  setProjectSteelOutOfWorkshop(input: $input)
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"projectId": 1,
			"identifierList": []string{
				"8",
			},
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取订单型钢详情集成测试--手持机
 */
func testProjectAdminRoleGetOrderSteelDetail(t *testing.T) {
	q := `
	query ($input: GetOrderSteelDetailInput!){
	  getOrderSteelDetail(input: $input) {
		id
		steel {id   
		  # 规格
		  specifcation{
			specification
		  }
		  # 型钢编码
		  code
		}
		#位置编码
		locationCode
	  }
	}
	`
	v := map[string]interface{} {
		"input": map[string]interface{}{
			"identifier": "8",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取未读消息总量集成测试--手持机
 */
func TestProjectAdminRoleGetMsgUnReadeTotal(t *testing.T) {
	q := `
		query {
		  getMsgUnReadeTotal # 未读消息总量
		}
	`
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}


/**
 * 项目管理员标记消息已读集成测试--手持机
 */
func testProjectAdminRoleGetMsgUnReadeTotal(t *testing.T) {
	me, _ := GetUserByToken(projectAdminTestCtx.DeviceToken)
	msgItem := msg.Msg{}
	err := model.DB.Model(&msgItem).Where("uid = ?", me.Id).
		First(&msgItem).
		Error
	assert.NoError(t, err)
	q := `
		mutation ($input: SetMsgReadedInput!){
		  setMsgBeRead(input: $input)
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{} {
			"idList": []int64 {msgItem.Id},
		},
	}
	_, err  = graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}
/**
 * 维修管理员获取日志列表集成测试
 */
func TestProjectAdminRoleGetLogList(t *testing.T) {
	q := `
		query ($input: GetLogListInput!){
		  getLogList(input: $input)  {
			list {
			  id
			  typeInfo {
				flag # 类型标志
				desc # 类型说明
	          }
			  content # 操作内容 
			  user {
				id 
				name # 操作用户
			  }
			  createdAt # 添加时间
			}
			total # 数量 
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"isShowAll": false,
			"page": 1,
			"pageSize": 10,
			"type": "DELETE",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取日类型志列表集成测试
 */
func TestProjectAdminRoleGetLogTypeList(t *testing.T) {
	q := `
		query {
		  getLogTypeList {
			desc # 类型说明 
			flag   #类型标志
		  }
		}
	`
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	assert.NoError(t, err)
}
/**
 * 项目管理员获取型钢详情列表集成测试
 */
func TestProjectAdminRoleGetProjectDetail(t *testing.T) {
	q := `
	query ($input: GetProjectDetailInput!){
	  getProjectDetail(input: $input){
		list{
		  id # 序号
		  orderSpecification {
			order {
			  id # 订单号
			}
		  }
		  steel {
			repository {
			  id 
			  name # 出货仓库
			}
			code # 型钢编码
			specifcation {
			  specification # 规格尺寸
			}
		  }
		  stateInfo {
			state 
			desc #状态说明 
		  }
		  outRepositoryAt # 出库时间
		  enterRepositoryAt #归库时间
		  enterWorkshopAt # 入场时间
		  outWorkshopAt # 出场时间
		  locationCode # 安装码
		  installationAt # 安装时间
		  useDays # 使用天数
		}
		total
	  }
	} 
	`
	v = map[string]interface{} {
		"input": map[string]interface{}{
			"isShowAll": false,
			"page": 1,
			"pageSize": 10,
			"state": 100,
			"outOfRepositoryAt": "2021-07-08T16:14:33+08:00",
			"locationCode": "2",
			"projectId": 2,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}
/**
 * 项目管理员获取型钢详情列表集成测试
 */
func TestProjectAdminRoleGetMaintenanceDetail(t *testing.T) {
	q := `
		query ($input: GetMaintenanceDetailInput!){
		  getMaintenanceDetail(input: $input) {
			list {
			  id
			  stateInfo {
				desc # 维修状态
				state
			  }
			  steel {
				repository{
				  id 
				  name #仓库
				}
				code # 型钢编码
				specifcation {
				  id
				  specification # 尺寸          
				}
				stateInfo {
				  desc # 当前状态说明
				  state # 当前状态
				}
			  }
			  enteredAt # 入厂时间
			  outedAt # 出厂时间
			  useDays # 维修天数
			}
			total #  数量 
			weight # 重量
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"isShowAll": false,
			"code": "GMSP-SJS01-00000",
			"page": 1,
			"pageSize": 2,
			"outMaintenanceAt": "2021-07-09T16:54:45+08:00",
			"state": 303,
			"repositoryId": 1,
			"specificationId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取维修的状态列表集成测试
 */
func TestProjectAdminRoleGetStateForMaintenance(t *testing.T) {
	q := `
		query {
		  getStateForMaintenance {
			state
			desc
		  }
		}
	`
	_, err := graphReqClient(q, v, roles.RoleMaintenanceAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员获取订单详情列表集成测试
 */
func TestProjectAdminRoleGetOrderDetailForBackEnd(t *testing.T) {
	q := `
		query ($input: GetOrderDetailForBackEndInput!) {
		  getOrderDetailForBackEnd(input: $input){
			list {
			  id
			  order {
					  orderNo # 订单号
					  project {
						id
						name # 项目名
					  }
					  repository {
						id
						name # 仓库名
					  }
			   
			  }
			  specification # 规格
			  total # 数量 
			  weight # 重量 
			}
			total # 数量 
			weight # 重量 
		  }
		}
	`
	v = map[string]interface{} {
		"input": map[string]interface{}{
			"isShowAll": false,
			"page": 1,
			"pageSize": 12,
			"projectId": 1,
			"repositoryId": 1,
			"specificationId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}
/**
 * 项目管理员获取仓库详情集成测试
 */
func TestProjectAdminRoleGetRepositoryDetail(t *testing.T) {
	q := `
	query ($input: GetRepositoryDetailInput!){
	  getRepositoryDetail(input: $input) {
		id
		name # 仓库名
		total # 数量 
		weight # 重量
		fee # 费用
	  }
	}
	`
	v = map[string]interface{} {
		"input": map[string]interface{}{
			"repositoryId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}

/**
 * 项目管理员删除订单集成测试
 */
func testProjectAdminRoleDeleteOrder(t *testing.T) {
	q := `
		mutation ($input: DeleteOrderInput!){
		  deleteOrder(input: $input) 
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface{}{
			"id": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}


/**
 * 项目管理员修改订单集成测试
 */
func testProjectAdminRoleEditOrder(t *testing.T) {
	q := `
		mutation ($input: EditOrderInput!) {
		  editOrder(input: $input) {
			id
		  }
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface{}{
			"id": 1,
			"partList": "",
			"steelList": []map[string]interface{}{
				{
					"specificationId": 1,
					"total": 1,
				},
			},
			"expectedReturnAt": "2021-08-22T13:49:33+08:00",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin)
	assert.NoError(t, err)
}
