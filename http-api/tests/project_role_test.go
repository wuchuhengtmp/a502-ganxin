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
 *
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
	v := map[string]interface{} {
		"input": map[string]interface{} {
			"projectId": 1,
			"state": 200,
			"specificationId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}
