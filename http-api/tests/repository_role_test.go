/**
 * @Desc    仓库管理员角色集成测试
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/2
 * @Listen  MIT
 */
package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"http-api/app/models/codeinfo"
	"http-api/app/models/devices"
	"http-api/app/models/orders"
	"http-api/app/models/roles"
	"http-api/app/models/steel_logs"
	"http-api/pkg/model"
	"http-api/seeders"
	"math/rand"
	"testing"
	"time"
)

// 仓库管理员测试上下文
var repositoryAdminTestCtx = struct {
	Token       string
	DeviceToken string
	Username    string
	Password    string
	// 用于个性规格记录
	EditSpecificationId   int64
	DeleteSpecificationId int64
	// 用于编辑的材料商家id
	EditMaterialId int64
	// 用于删除材料商家
	DeleteMaterialId int64
	// 用于编辑制造商家
	EditManufacturerId int64
	// 用于删除制造商家测试
	DeleteManufacturerId int64
	// 用于编辑物流公司
	EditExpressId int64
	// 删除物流公司
	DeleteExpressId int64
}{
	Username: seeders.RepositoryAdmin.Username,
	Password: seeders.RepositoryAdmin.Password,
}

/**
 * 仓库管理员登录测试
 */
func TestRepositoryAdminRoleLogin(t *testing.T) {
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
		"phone":    repositoryAdminTestCtx.Username,
		"password": repositoryAdminTestCtx.Password,
	}
	res, err := graphReqClient(query, variables, roles.RoleRepositoryAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	repositoryAdminTestCtx.Token = tokenInfo["accessToken"].(string)
}

/**
 * 仓库管理员登录测试-手持机
 */
func TestRepositoryAdminRoleDeviceLogin(t *testing.T) {
	query := `
		mutation ($phone: String!, $password: String!, $mac: String!){
		  login (phone: $phone, password: $password, mac: $mac) {
			accessToken
		  }
		}
	`
	variables := map[string]interface{}{
		"phone":    repositoryAdminTestCtx.Username,
		"password": repositoryAdminTestCtx.Password,
		"mac":      "123:1242:1242:12412",
	}
	res, err := graphReqClient(query, variables, roles.RoleRepositoryAdmin)
	hasError(t, err)
	token := res["login"]
	tokenInfo := token.(map[string]interface{})
	repositoryAdminTestCtx.DeviceToken = tokenInfo["accessToken"].(string)
}

/**
 * 仓库管理员获取公司列表集成测试
 */
func TestRepositoryAdminRoleGetAllCompany(t *testing.T) {
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
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题"))
	}
}

/**
 * 仓库管理员获取公司人员列表集成测试
 */
func TestRepositoryAdminRoleGetCompanyUsers(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员获取仓库列表集成测试
 */
func TestRepositoryAdminRoleGetRepository(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员创建规格集成测试
 */
func TestRepositoryAdminRoleCreateSpecification(t *testing.T) {
	q := `
		mutation createSpecificationMutation($input: CreateSpecificationInput!) {
			createSpecification(input: $input) {
			id
			length
			weight
			type
			isDefault
			specification
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"length":    rand.Intn(100),
			"weight":    rand.Intn(100),
			"type":      "type_test_for_repositoryRole",
			"isDefault": false,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	data := res["createSpecification"].(map[string]interface{})
	repositoryAdminTestCtx.EditSpecificationId = int64(data["id"].(float64))
	repositoryAdminTestCtx.DeleteSpecificationId = int64(data["id"].(float64))
}

/**
 * 仓库管理员获取规格列表集成测试
 */
func TestRepositoryAdminRoleGetSpecification(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
* 仓库管理员修改规格集成测试
 */
func TestRepositoryAdminRoleEditSpecification(t *testing.T) {
	q := `
		mutation editSpecificationMutation($input: EditSpecificationInput !) {
			editSpecification(input: $input) {
				id
				isDefault
				specification
				weight
				length
				type
			}
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditSpecificationId,
			"weight":    rand.Intn(100),
			"length":    rand.Float64(),
			"type":      "test_for_repositoryRole",
			"isDefault": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员删除规格集成测试
 */
func TestRepositoryAdminRoleDeleteSpecification(t *testing.T) {
	q := `
		mutation deleteSpecification($id: Int!) {
			deleteSpecification(id: $id)
		}
	`
	v := map[string]interface{}{
		"id": repositoryAdminTestCtx.DeleteSpecificationId,
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员创建材料商集成测试
 */
func TestRepositoryAdminRoleCreateCodeInfo(t *testing.T) {
	q := `
		mutation createMaterialManufacturerMutation ($input: CreateMaterialManufacturerInput!){
		  createMaterialManufacturer(input: $input) {
			id
			name
			
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "name_test_for_repositoryRoleCreateInfoCode",
			"remark":    "remark_for_repositoryRoleCreateInfoTest",
			"isDefault": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "name_test_for_repositoryRoleCreateInfoCode",
			"remark":    "remark_for_repositoryRoleCreateInfoTest",
			"isDefault": false,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	data := res["createMaterialManufacturer"].(map[string]interface{})
	id := data["id"].(float64)
	repositoryAdminTestCtx.EditMaterialId = int64(id)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "name_test_for_repositoryRoleCreateInfoCode",
			"remark":    "remark_for_repositoryRoleCreateInfoTest",
			"isDefault": false,
		},
	}
	res, err = graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	data = res["createMaterialManufacturer"].(map[string]interface{})
	id = data["id"].(float64)
	repositoryAdminTestCtx.DeleteMaterialId = int64(id)
}

/**
 * 仓库管理员获取材料商列表集成测试
 */
func TestRepositoryAdminRoleGetMaterialManufacturers(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员编辑材料商集成测试
 */
func TestRepositoryAdminRoleEditMaterialManufacturers(t *testing.T) {
	q := `mutation editMaterialManufacturerMutation($input: EditMaterialManufacturerInput!) {
		  editMaterialManufacturer(input: $input){
			id
			name
			remark
			isDefault
		  }
		}`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditMaterialId,
			"name":      "name_for_repositoryRoleTest",
			"remark":    "remark_for_repositoryRoleTest",
			"isDefault": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditMaterialId,
			"name":      "name_for_repositoryRoleTest",
			"remark":    "remark_for_repositoryRoleTest1",
			"isDefault": false,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员删除材料商集成测试
 */
func TestRepositoryAdminRoleDeleteMaterialManufacturers(t *testing.T) {
	q := `
		mutation deleteMaterialManufacturer($deleteId: Int!){
		  deleteMaterialManufacturer(id: $deleteId)
		}
	`
	v := map[string]interface{}{
		"deleteId": repositoryAdminTestCtx.DeleteMaterialId,
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	hasError(t, err)
}

/**
 * 仓库管理员创建制造商集成测试
 */
func TestCompanyRepositoryRoleCreateManufacturer(t *testing.T) {
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	var cs []*codeinfo.CodeInfo
	model.DB.Model(&codeinfo.CodeInfo{}).
		Where("type = ? AND company_id = ?", codeinfo.Manufacturer, me.CompanyId).
		Find(&cs)
	q := `
		mutation createManufacturerMutation($input: CreateManufacturerInput!) {
		  createManufacturer(input: $input) {
			id
			name
			isDefault
		  }
		}
	`
	name := fmt.Sprintf("name_for_createManufacturerByRepositoryRole_%s", fmt.Sprintf("%d", time.Now().UnixNano()))
	remark := "remark_for_createManufacturerTest"
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal(err.Error())
	}
	newCodeInfo := codeinfo.CodeInfo{}
	err = model.DB.Model(&codeinfo.CodeInfo{}).
		Where("name = ?", name).First(&newCodeInfo).
		Error
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, name, newCodeInfo.Name)
	assert.Equal(t, remark, newCodeInfo.Remark)
	if len(cs) > 0 && newCodeInfo.IsDefault != isDefault {
		assert.Equal(t, !isDefault, newCodeInfo.IsDefault)
	}
	if len(cs) == 0 {
		assert.Equal(t, true, newCodeInfo.IsDefault)
	}
	if len(cs) >= 2 {
		model.DB.Model(&codeinfo.CodeInfo{}).
			Where("company_id = ? AND type = ? AND is_default = ?", me.CompanyId, codeinfo.Manufacturer, true).
			Find(&cs)
		assert.Len(t, cs, 1)
	}
	// 保存id用于编辑制造商测试
	data := res["createManufacturer"].(map[string]interface{})
	id := data["id"].(float64)
	repositoryAdminTestCtx.EditManufacturerId = int64(id)
	repositoryAdminTestCtx.DeleteManufacturerId = int64(id)
}

/**
 * 仓库管理员获取制造商集成测试
 */
func TestCompanyRepositoryRoleGetManufacturer(t *testing.T) {
	q := `query getManufacturersQuery{
		getManufacturers{
		id
		isDefault
		name
		remark
	  }
	}`
	v := map[string]interface{}{}
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员获取制造商集成测试")
	}
	assertCompanyIdForGetManufacturers(t, res, repositoryAdminTestCtx.Token)
}

/**
 * 仓库管理员编辑制造商集成测试
 */
func TestRepositoryAdminRoleEditManufacturers(t *testing.T) {
	q := `
		mutation editManufacturerMutation($input: EditManufacturerInput! ){
		  editManufacturer(input: $input) {
			id
			name
			isDefault
			remark
		  }
		}
	`
	name := "name_for_repositoryRoleTest"
	remark := "remark_form_repositoryRoleTest"
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditManufacturerId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:公司仓库管理员编辑制造商集成测试")
	}
	c := codeinfo.CodeInfo{
		ID: repositoryAdminTestCtx.EditManufacturerId,
	}
	_ = c.GetSelf()
	assert.Equal(t, name, c.Name)
	assert.Equal(t, isDefault, c.IsDefault)
	assert.Equal(t, remark, c.Remark)
	// 只能有一个默认选项
	var cs []*codeinfo.CodeInfo
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	model.DB.Model(&codeinfo.CodeInfo{}).
		Where("type = ? AND company_id = ? AND is_default = ?", codeinfo.Manufacturer, me.CompanyId, true).
		Find(&cs)
	assert.Len(t, cs, 1)
	c = *cs[0]
	assert.Equal(t, name, c.Name)
	assert.Equal(t, isDefault, c.IsDefault)
	assert.Equal(t, remark, c.Remark)
}

/**
 * 仓库管理员删除制造商集成测试
 */
func TestRepositoyAdminRoleDeleteManufacturers(t *testing.T) {
	q := `
		mutation deleteManufacturerMutation($id: Int!) {
		  deleteManufacturer(id: $id) 
		}
	`
	v := map[string]interface{}{
		"id": repositoryAdminTestCtx.DeleteManufacturerId,
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员删除制造商集成测试")
	}
	// 断言没有这条数据了
	var cs []*codeinfo.CodeInfo
	model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", repositoryAdminTestCtx.DeleteManufacturerId).Find(&cs)
	assert.Len(t, cs, 0)
	// 断言有新的默认制造商家了
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	model.DB.Model(&codeinfo.CodeInfo{}).Where("company_id = ? AND type = ?", me.CompanyId, codeinfo.Manufacturer).Find(&cs)
	if len(cs) > 0 {
		c := codeinfo.CodeInfo{}
		err := model.DB.
			Model(&codeinfo.CodeInfo{}).
			Where("company_id = ? AND type = ? AND is_default = ?", me.CompanyId, codeinfo.Manufacturer, true).
			First(&c).
			Error
		assert.NoError(t, err)
	}
}

/**
 * 仓库管理员创建物流商集成测试
 */
func TestRepositoryAdminRoleCreateExpress(t *testing.T) {
	q := `
		mutation createExpressMutation($input: CreateExpressInput!) {
		  createExpress(input: $input) {
			id
			name
			remark
			isDefault
		  }
		}
	`
	// 默认断言
	name := "name_for_repositoryCreateExpress_" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark := "remark_for_repositoryCreateExpress"
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	c := codeinfo.CodeInfo{}
	if err := model.DB.Model(&codeinfo.CodeInfo{}).Where("name = ?", name).First(&c).Error; err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	assert.Equal(t, c.Name, name)
	assert.Equal(t, c.IsDefault, isDefault)
	assert.Equal(t, c.Type, codeinfo.ExpressCompany)
	assert.Equal(t, c.CompanyId, me.CompanyId)
	var cs []*codeinfo.CodeInfo
	model.DB.Model(&codeinfo.CodeInfo{}).
		Where("company_id = ? AND type = ? AND is_default = ?", me.CompanyId, codeinfo.ExpressCompany, true).
		Find(&cs)
	assert.Len(t, cs, 1)
	assert.Equal(t, cs[0].ID, c.ID)
	// 非默认断言
	name = "name_for_repositoryCreateExpress_" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark = "remark_for_repositoryCreateExpress"
	isDefault = false
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	c = codeinfo.CodeInfo{}
	if err := model.DB.Model(&codeinfo.CodeInfo{}).Where("name = ?", name).First(&c).Error; err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	assert.Equal(t, c.Name, name)
	assert.Equal(t, c.IsDefault, isDefault)
	assert.Equal(t, c.Type, codeinfo.ExpressCompany)
	assert.Equal(t, c.CompanyId, me.CompanyId)
	model.DB.Model(&codeinfo.CodeInfo{}).
		Where("company_id = ? AND type = ? AND is_default = ?", me.CompanyId, codeinfo.ExpressCompany, true).
		Find(&cs)
	assert.Len(t, cs, 1)
	data := res["createExpress"].(map[string]interface{})
	id := data["id"].(float64)
	repositoryAdminTestCtx.EditExpressId = int64(id)
	repositoryAdminTestCtx.DeleteExpressId = int64(id)
}

/**
 * 仓库管理员获取物流列表集成测试
 */
func TestRepositoryAdminRoleGetExpressList(t *testing.T) {
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
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员获取物流列表集成测试")
	}
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
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
 * 仓库管理员编辑物流集成测试
 */
func TestRepositoryAdminRoleEditExpress(t *testing.T) {
	q := `
		mutation editExpressMutation($input: EditExpressInput!) {
		  editExpress(input: $input) {
			id
			isDefault
			remark
			name
		  }
		}
	`
	name := "name_for_companyRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark := "remark_for_companyRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditExpressId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员编辑物流集成测试")
	}
	record := codeinfo.CodeInfo{}
	err = model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", repositoryAdminTestCtx.EditExpressId).First(&record).Error
	if err != nil {
		t.Fatal("failed:仓库管理员编辑物流集成测试")
	}
	assert.Equal(t, record.IsDefault, isDefault)
	assert.Equal(t, record.Name, name)
	assert.Equal(t, record.Type, codeinfo.ExpressCompany)
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	assert.Equal(t, record.CompanyId, me.CompanyId)
	assert.Equal(t, remark, record.Remark)
	// 不是默认选项断言
	name = "name_for_repositoryRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark = "remark_for_repositoryRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	isDefault = false
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"id":        repositoryAdminTestCtx.EditExpressId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员编辑物流集成测试")
	}
	var cs []codeinfo.CodeInfo

	model.
		DB.
		Model(&codeinfo.CodeInfo{}).
		Where("company_id = ? AND type = ?", me.CompanyId, codeinfo.ExpressCompany).
		Find(&cs)
	if len(cs) > 1 {
		model.
			DB.
			Model(&codeinfo.CodeInfo{}).
			Where("company_id = ? AND is_default = ? AND type = ?", me.CompanyId, true, codeinfo.ExpressCompany).
			Find(&cs)
		assert.Len(t, cs, 1)
	}
}

func TestRepositoryAdminRoleDeleteExpress(t *testing.T) {
	q := `
		mutation deleteExpressMutation($id: Int!){
			deleteExpress(id: $id)
		}
	`
	v := map[string]interface{}{
		"id": repositoryAdminTestCtx.DeleteExpressId,
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed: 仓库管理员编辑物流集成测试")
	}
	// 断言已删除
	c := codeinfo.CodeInfo{}
	err = model.DB.Model(&c).Where("id = ?", repositoryAdminTestCtx.DeleteExpressId).First(&c).Error
	assert.Error(t, err)
	var cs []codeinfo.CodeInfo
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
	model.DB.Model(&codeinfo.CodeInfo{}).
		Where("company_id = ? AND type = ?", me.CompanyId, codeinfo.ExpressCompany).
		Find(&cs)
	if len(cs) > 0 {
		// 断言有指定新的默认
		err := model.DB.Model(&codeinfo.CodeInfo{}).
			Where("company_id = ? AND type = ? AND is_default = ?", me.CompanyId, codeinfo.ExpressCompany, true).
			First(&c).
			Error
		assert.NoError(t, err)
	}
}

/**
 * 仓库管理员获取价格集成测试
 */
func TestRepositoryAdminRoleGetPrice(t *testing.T) {
	q := `
		 query getPriceQuery {
		  getPrice
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员获取集成测试")
	}
}

/**
 * 仓库管理员设备登录集成测试
 */
func TestRepositoryAdminRoleLoginDevice(t *testing.T) {
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
		"phone":    "13427969606",
		"password": "12345678",
		"mac":      "123:1242:1242:12412",
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员设备登录集成测试")
	}
}

/**
 * 仓库管理员获取设备列表集成测试
 */
func TestRepositoryAdminGetDeviceList(t *testing.T) {
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
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("failed:仓库管理员获取设备列表集成测试")
	}
	// 断言响应的数据就是用户的仓库名下的
	me, _ := GetUserByToken(repositoryAdminTestCtx.Token)
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
 * 仓库管理员入库型钢集成测试
 */
func TestRepositoryAdminCreateSteel(t *testing.T) {
	var oldTotalLogs int64
	model.DB.Model(&steel_logs.SteelLog{}).Count(&oldTotalLogs)
	q := `
		mutation createSteelMutation($input: CreateSteelInput! ){
		  createSteel(input: $input) {
			id
			state
			specifcation{
			  id
			  specification
			}
			turnover
			producedDate
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"identifierList": []interface{}{
				fmt.Sprintf("%d", time.Now().UnixNano()),
				fmt.Sprintf("%d", time.Now().UnixNano()),
				fmt.Sprintf("%d", time.Now().UnixNano()),
			},
			"repositoryId":           1,
			"specificationId":        1,
			"manufacturerId":         4,
			"materialManufacturerId": 1,
			"producedDate":           "2021-06-11T10:08:42+08:00",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal(err)
	}
	// 断言型钢操作日志增量3
	var newTotalLogs int64
	model.DB.Model(&steel_logs.SteelLog{}).Count(&newTotalLogs)
	assert.Equal(t, oldTotalLogs+3, newTotalLogs)
}

/**
 * 仓库管理员获取型钢列表集成测试
 */
func TestRepositoryAdminGetSteelList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("仓库管理员获取型钢列表集成测试")
	}

}

/**
 * 仓库管理员设置密码集成测试
 */
func TestRepositoryAdminSetPassword(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("仓库管理员设置密码集成测试")
	}
}

/**
 * 仓库管理员获取我的信息集成测试
 */
func TestRepositoryAdminGetMyInfo(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	if err != nil {
		t.Fatal("仓库管理员获取我的信息集成测试")
	}
}

/**
 * 仓库管理员获取项目列表集成测试
 */
func TestRepositoryAdminGetProjectList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	assert.NoError(t, err)
}

/**
 * 仓库管理员获取订单列表集成测试
 */
func TestRepositoryAdminGetOrderList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin)
	assert.NoError(t, err)
}

/**
 * 仓库管理员获取订单列表集成测试-手持机
 */
func TestRepositoryAdminDeviceGetOrderList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理员获取订单详情集成测试-手持机
 */
func TestRepositoryAdminDeviceGetOrderDetail(t *testing.T) {
	me, err := GetUserByToken(repositoryAdminTestCtx.DeviceToken)
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
	_, err = graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理员确认或拒绝订单集成测试-手持设备
 */
func TestRepositoryAdminConfirmOrRejectOrder(t *testing.T) {
	me, err := GetUserByToken(repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
	o := orders.Order{}
	err = model.DB.Model(&o).Where("company_id = ? AND state = ?", me.CompanyId, orders.StateToBeConfirmed).First(&o).Error
	assert.NoError(t, err)
	q := `
		mutation ($input: ConfirmOrderInput!){
		  confirmOrRejectOrder(input: $input) {
			id
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":       o.Id,
			"isAccess": true,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
	t.Run("仓库管理员获取待出库详细信息-手持设备", testRepositoryAdminGetProjectOrder2WorkshopDetail)
	t.Run("仓库管理员型钢出库集成测试-手持机", testRepositoryAdminRoleSetProjectOrder2Workshop)
	t.Run("项目管理员型钢入场集成测试-手持设备", testProjectAdminRoleSetProjectEnterWorkshop)
	t.Run("项目管理员获取项目规格列表集成测试--手机", testProjectAdminRoleGetProjectSpecificationDetail)
	t.Run("项目管理员获取项目详情列表集成测试--手机", testProjectAdminRoleGetProjectSteelDetail)
	t.Run("项目管理员获取项目型钢状态列表集成测试", testProjectAdminRoleGetProjectSteelStateList)
	t.Run("项目管理员获取最大安装码表集成测试--手机", testProjectAdminRoleGetMaxLocationCode)
	t.Run("项目管理员安装码是否可用集成测试--手持机", testProjectAdminRoleIsAccessLocationCode)
	t.Run("项目管理员安装型钢集成测试-手持机", testProjectAdminRoleInstallSteel)
	t.Run("项目管理员待修改型钢详细信息集成测试-手持机", testProjectAdminRoleGetProjectSteel2BeChangeDetail)
	t.Run("项目管理员 项目中的型钢查询 集成测试-手持机", testProjectAdminRoleGetProjectSteel2BeChange)
	t.Run("项目管理员 修改项目型钢状态 集成测试-手持机", testProjectAdminRoleSetProjectSteelState)
	t.Run("获取可出场的项目列表", testProjectAdminRoleGetOutOfWorkshopProjectList)
	t.Run("项目管理员获取可出场的项目列表集成测试--手持机", testProjectAdminRoleSetProjectSteelOutOfWorkshop)
	t.Run("项目管理员获取订单型钢详情集成测试--手持机", testProjectAdminRoleGetOrderSteelDetail)
	t.Run("仓库管理中获取可归库的项目列表集成测试--手持机", testRepositoryAdminRoleGetEnterRepositoryProjectList)
	t.Run("项目管理员获取修改项目状态列表集成测试--手持机", testProjectAdminRoleGetProjectStateListForChanged)
}

/**
 * 仓库管理员获取待出库详细信息-手持设备
 */
func testRepositoryAdminGetProjectOrder2WorkshopDetail(t *testing.T) {
	q := `
		query ($input: ProjectOrder2WorkshopDetailInput!){
		   getProjectOrder2WorkshopDetail(input: $input) {
			#  这是列表
			list{
				id
				identifier
				# 这是规格相关信息
				specifcation{
				  specification
				  id
				}
				state
				turnover

			  }
			# 数量
			total
			# 重量
			totalWeight 
			}
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface{}{
			"orderId": 1,
			"identifierList": []string {
				"11",
			},
		},
	}
	_ ,err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理员获取获取可以出库的订单列表集成测试-手持机
 */
func TestRepositoryAdminGetTobeSendWorkshopOrderList(t *testing.T) {
	q := `
		query {
		   getTobeSendWorkshopOrderList{
			id
			state
		  }
		}
	`
	v := map[string]interface{} {}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理员型钢快速查询集成测试-手持机
 */
func TestRepositoryAdminGetSteelDetail(t *testing.T)  {
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
	v := map[string]interface{} {
		"input": map[string]interface{} {
			"identifier": "8",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理员型钢出库集成测试-手持机
 */
func testRepositoryAdminRoleSetProjectOrder2Workshop(t *testing.T) {
	q := `
		mutation ($input: ProjectOrder2WorkshopInput!) {
		  setProjectOrder2Workshop(input: $input) {
			id
		  }
		}
	`
	v := map[string]interface{} {
		"input": map[string]interface{} {
			"expressCompanyId": 12,
			"expressNo": "124333123",
			"identifierList": []string{
				"8",
			},
			"orderId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库多个快速查询集成测试-手持机
 */
func TestRepositoryAdminRoleGetMultipleSteelDetail(t *testing.T) {
	q := `
		query ($input: GetMultipleSteelDetailInput!){
		  getMultipleSteelDetail(input: $input) {
			 id
			code
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
	v := map[string]interface{} {
		"input": map[string]interface{} {
			"identifierList": []string {
				"8",
				"9",
			},
		},
	}
	_, err := graphReqClient(q, v, roles.RoleProjectAdmin, projectAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 *仓库管理员获取消息列表集成测试-手持机
 */
func TestRepositoryAdminRoleGetMsgList(t *testing.T) {
	q := `query {
		  getMsgList{
			id
			isRead
			content
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}


/**
 * 仓库管理中获取可归库的项目列表集成测试--手持机
 */
func testRepositoryAdminRoleGetEnterRepositoryProjectList(t *testing.T) {
	q := `
		query {
		  getEnterRepositoryProjectList{
			id
			name
		  }
		}
	`
	v := map[string]interface{}{ }
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}

/**
 * 仓库管理中获取可归库的项目列表集成测试--手持机
 */
func testRepositoryAdminRole(t *testing.T) {
	q := `
		query ($input: GetEnterRepositorySteelDetailInput!){
		  getEnterRepositorySteelDetail(input: $input) {
			id
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{} {
			"identifier": "8",
			"projectId": 1,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleRepositoryAdmin, repositoryAdminTestCtx.DeviceToken)
	assert.NoError(t, err)
}
