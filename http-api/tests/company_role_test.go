/**
 * @Desc    公司管理员角色集成测试
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
	"http-api/app/models/codeinfo"
	"http-api/app/models/configs"
	"http-api/app/models/devices"
	"http-api/app/models/logs"
	"http-api/app/models/roles"
	"http-api/pkg/model"
	"http-api/seeders"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// 超级管理员测试上下文
var companyAdminTestCtx = struct {
	// token 用于角色鉴权
	Token string
	// 用于删除的公司id
	DeleteCompanyId int64
	// 账号
	Username string
	// 密码
	Password string
	// 用于删除的公司员工id
	DeleteCompanyUserId int64
	// 用于编辑的公司员工id
	EditCompanyUserId int64
	// 用于删除仓库测试
	DeleteRepositoryId int64
	// 用于编辑规格
	EditSpecificationId int64
	// 用于删除规格
	DeleteSpecificationId int64
	// 用于编辑的材料商id
	EditMaterialId int64
	// 用于删除的材料商id
	DeleteMaterialId int64
	// 用于编辑制造商ID
	EditManufacturerId int64
	// 用于删除制造商家ID
	DeleteManufacturerId int64
	// 用于编辑物流公司ID
	EditExpressId int64
	// 删除物流公司ID
	DeleteExpressId int64
	// 编辑设备id
	EditDeviceId int64
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
	variables := map[string]interface{}{
		"phone":    companyAdminTestCtx.Username,
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
func TestCompanyAdminRoleGetAllCompany(t *testing.T) {
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
	hasError(t, err)
	if len(res) != 1 {
		hasError(t, fmt.Errorf("期望返回一条公司数据，结果不是，要么是没有数据， 要么是数据权限作用域限制出了问题"))
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
		"input": map[string]interface{}{
			"id":                2,
			"name":              "name_for_companyRoleEditTest",
			"pinYin":            "GMSP",
			"symbol":            "企业宗旨修改测试_companyRoleEditTest",
			"logoFileId":        1,
			"backgroundFileId":  2,
			"isAble":            true,
			"phone":             seeders.CompanyAdmin.Username,
			"wechat":            "12345678",
			"startedAt":         "2021-12-31 00:00:00",
			"endedAt":           "2022-12-31 00:00:00",
			"adminName":         "username_change_test_with_company_role",
			"adminPassword":     seeders.CompanyAdmin.Password,
			"adminAvatarFileId": 4,
			"adminPhone":        seeders.CompanyAdmin.Username,
			"adminWechat":       "admin_wechat_change_test_with_company_role",
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 添加公司人员集成测试
 */
func TestCompanyAdminRoleCreateCompanyUser(t *testing.T) {
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
		"input": map[string]interface{}{
			"name":     "username _for_TesCreateCompanyUser",
			"phone":    fmt.Sprintf("1342%s", fmt.Sprintf("%d", time.Now().UnixNano())[8:15]),
			"password": "12345678",
			"avatarId": 1,
			"role":     "repositoryAdmin",
			"wechat":   "wechat_for_testCreateCompanyUser",
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
	user := res["createCompanyUser"].(map[string]interface{})
	id := user["id"].(float64)
	companyAdminTestCtx.DeleteCompanyUserId = int64(id)
	companyAdminTestCtx.EditCompanyUserId = int64(id)
}

/**
 * 获取公司人员列表集成测试
 */
func TestCompanyAdminRoleGetCompanyUsers(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员修改公司人员集成测试
 */
func TestCompanyAdminRoleEditCompanyUser(t *testing.T) {
	q := `mutation editCompanyUserMutaion ($input: EditCompanyUserInput!){
			editCompanyUser(input: $input) {
			id
			role {
			  id
				name
			  tag
			}
			phone
			wechat
			avatar {
			  id
			  url
			}
			isAble
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     companyAdminTestCtx.EditCompanyUserId,
			"name":   "change_name_for_editCompanyUser",
			"phone":  fmt.Sprintf("1342%s", fmt.Sprintf("%d", time.Now().UnixNano())[8:15]),
			"roleId": 2,
			"isAble": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员删除公司人员集成测试
 */
func TestCompanyAdminRoleDeleteCompanyUser(t *testing.T) {
	q := `mutation deleteCompanyUserMutation($uid: Int!){
		  deleteCompanyUser(uid: $uid)
		}
	`
	v := map[string]interface{}{
		"uid": companyAdminTestCtx.DeleteCompanyUserId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 *  公司管理员添加仓库集成测试
 */
func TestCompanyAdminRoleCreateRepository(t *testing.T) {
	q := `
		mutation createRepository($input: CreateRepositoryInput!) {
		  createRepository(input: $input) {
			id
			weight
			pinYin
			address
			total
			weight
			remark
			isAble
			total
			leaders {
				id
				name
				wechat
				phone
			}
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":              "reposistory_name_for_test",
			"remark":            "",
			"address":           "address_for_createAddress",
			"repositoryAdminId": 3,
			"pinYin":            "pintYin_for_createTest",
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
	r := res["createRepository"].(map[string]interface{})
	id := r["id"].(float64)
	companyAdminTestCtx.DeleteRepositoryId = int64(id)
}

/**
 * 公司管理员获取获取仓库列表集成测试
 */
func TestCompanyAdminRoleGetRepository(t *testing.T) {
	q := `
		 query {
		  getRepositoryList {
			id
			isAble
			weight
			leaders{
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员删除仓库集成测试
 */
func TestCompanyAdminRoleDeleteRepository(t *testing.T) {
	q := `
		mutation deleteRepositoryMutation ($id: Int!) {
		  deleteRepository(repositoryId: $id)
		}
	`
	v := map[string]interface{}{
		"id": companyAdminTestCtx.DeleteRepositoryId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}

/**
 * 公司管理员添加规格集成测试
 */
func TestCompanyAdminRoleCreateSpecification(t *testing.T) {
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
			"type":      "type_test_for_companyAdminRole",
			"isDefault": false,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
	data := res["createSpecification"].(map[string]interface{})
	companyAdminTestCtx.EditSpecificationId = int64(data["id"].(float64))
	companyAdminTestCtx.DeleteRepositoryId = int64(data["id"].(float64))
}

/**
 * 公司管理员获取规格集成测试
 */
func TestCompanyAdminRoleGetSpecification(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员修改规格集成测试
 */
func TestCompanyAdminRoleEditSpecification(t *testing.T) {
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
			"id":        companyAdminTestCtx.EditSpecificationId,
			"weight":    rand.Intn(100),
			"length":    rand.Float64(),
			"type":      "test_for_CompanyRole",
			"isDefault": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员删除规格集成测试
 */
func TestCompanyAdminRoleDeleteSpecification(t *testing.T) {
	q := `
		mutation deleteSpecification($id: Int!) {
			deleteSpecification(id: $id)
		}
	`
	v := map[string]interface{}{
		"id": companyAdminTestCtx.DeleteRepositoryId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员添加材料商集成测试
 */
func TestCompanyAdminRoleCreatCodeInfo(t *testing.T) {
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
			"name":      "name_test_for_companyRoleCreateInfoCode",
			"remark":    "remark_for_companyRoleCreateInfoTest",
			"isDefault": true,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "name_test_for_companyRoleCreateInfoCode",
			"remark":    "remark_for_companyRoleCreateInfoTest",
			"isDefault": false,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
	body := res["createMaterialManufacturer"].(map[string]interface{})
	id := body["id"].(float64)
	companyAdminTestCtx.EditMaterialId = int64(id)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"name":      "name_test_for_companyRoleCreateInfoCode",
			"remark":    "remark_for_companyRoleCreateInfoTest",
			"isDefault": false,
		},
	}
	res, err = graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
	body = res["createMaterialManufacturer"].(map[string]interface{})
	id = body["id"].(float64)
	companyAdminTestCtx.DeleteMaterialId = int64(id)
}

/**
 * 公司管理员获取材料商列表集成测试
 */
func TestCompanyAdminRoleGetMaterialManufacturers(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员删除材料商集成测试
 */
func TestCompanyAdminRoleDeleteMaterialManufacturers(t *testing.T) {
	q := `mutation deleteMaterialManufacturer($deleteId: Int!){
		  deleteMaterialManufacturer(id: $deleteId)
		}
	`
	v := map[string]interface{}{
		"deleteId": companyAdminTestCtx.DeleteMaterialId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	hasError(t, err)
}

/**
 * 公司管理员添加制造商集成测试
 */
func TestCompanyAdminRoleCreateManufacturer(t *testing.T) {
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
	name := fmt.Sprintf("name_for_createManufacturerTest_%s", fmt.Sprintf("%d", time.Now().UnixNano()))
	remark := "remark_for_createManufacturerByCompanyRoleTest"
	isDefault := false
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      name,
			"remark":    remark,
			"isDefault": rand.Intn(2) == 1,
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
	data := res["createManufacturer"].(map[string]interface{})
	id := data["id"].(float64)
	companyAdminTestCtx.EditManufacturerId = int64(id)
	companyAdminTestCtx.DeleteManufacturerId = int64(id)
}

/**
 * 公司管理员获取制造商列表集成测试
 */
func TestCompanyAdminRoleGetManufacturers(t *testing.T) {
	q := `query getManufacturersQuery{
		getManufacturers{
		id
		isDefault
		name
		remark
	  }
	}`
	v := map[string]interface{}{}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员获取制造商列表集成测试")
	}
	assertCompanyIdForGetManufacturers(t, res, companyAdminTestCtx.Token)
}

/**
 * 公司管理员编辑制造商集成测试
 */
func TestCompanyAdminRoleEditManufacturers(t *testing.T) {
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
	name := "name_for_companyRoleTest"
	remark := "remark_form_CompanyRoleTest"
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":        companyAdminTestCtx.EditManufacturerId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员编辑制造商集成测试")
	}
	c := codeinfo.CodeInfo{
		ID: companyAdminTestCtx.EditManufacturerId,
	}
	_ = c.GetSelf()
	assert.Equal(t, name, c.Name)
	assert.Equal(t, isDefault, c.IsDefault)
	assert.Equal(t, remark, c.Remark)
	// 只能有一个默认选项
	var cs []*codeinfo.CodeInfo
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
 * 公司管理员删除制造商集成测试
 */
func TestCompanyAdminRoleDeleteManufacturers(t *testing.T) {
	q := `
		mutation deleteManufacturerMutation($id: Int!) {
		  deleteManufacturer(id: $id) 
		}
	`
	v := map[string]interface{}{
		"id": companyAdminTestCtx.DeleteManufacturerId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员删除制造商集成测试")
	}
	// 断言没有这条数据了
	var cs []*codeinfo.CodeInfo
	model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", companyAdminTestCtx.DeleteManufacturerId).Find(&cs)
	assert.Len(t, cs, 0)
	// 断言有新的默认制造商家了
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
 * 公司管理员创建物流商集成测试
 */
func TestCompanyAdminRoleCreateExpress(t *testing.T) {
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
	name := "name_for_companyCreateExpress_" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark := "remark_for_companyCreateExpress"
	isDefault := true
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	c := codeinfo.CodeInfo{}
	if err := model.DB.Model(&codeinfo.CodeInfo{}).Where("name = ?", name).First(&c).Error; err != nil {
		t.Fatal("failed:创建物流公司集成测试失败")
	}
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
	name = "name_for_companyCreateExpress_" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark = "remark_for_companyCreateExpress"
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
	companyAdminTestCtx.EditExpressId = int64(id)
	companyAdminTestCtx.DeleteExpressId = int64(id)
}

/**
 * 公司管理员获取物流列表集成测试
 */
func TestCompanyAdminRoleGetExpressList(t *testing.T) {
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
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员获取物流列表集成测试")
	}
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
 * 公司管理员编辑物流集成测试
 */
func TestCompanyAdminRoleEditExpress(t *testing.T) {
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
			"id":        companyAdminTestCtx.EditExpressId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员编辑物流集成测试")
	}
	record := codeinfo.CodeInfo{}
	err = model.DB.Model(&codeinfo.CodeInfo{}).Where("id = ?", companyAdminTestCtx.EditExpressId).First(&record).Error
	if err != nil {
		t.Fatal("failed:公司管理员编辑物流集成测试")
	}
	assert.Equal(t, record.IsDefault, isDefault)
	assert.Equal(t, record.Name, name)
	assert.Equal(t, record.Type, codeinfo.ExpressCompany)
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
	assert.Equal(t, record.CompanyId, me.CompanyId)
	assert.Equal(t, remark, record.Remark)
	// 不是默认选项断言
	name = "name_for_companyRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	remark = "remark_for_companyRoleEditTest" + fmt.Sprintf("%d", time.Now().UnixNano())
	isDefault = false
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"id":        companyAdminTestCtx.EditExpressId,
			"name":      name,
			"remark":    remark,
			"isDefault": isDefault,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员编辑物流集成测试")
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

/**
 * 公司管理员编辑物流集成测试
 */
func TestCompanyAdminRoleDeleteExpress(t *testing.T) {
	q := `
		mutation deleteExpressMutation($id: Int!){
			deleteExpress(id: $id)
		}
	`
	v := map[string]interface{}{
		"id": companyAdminTestCtx.DeleteExpressId,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed: 公司管理员编辑物流集成测试")
	}
	// 断言已删除
	c := codeinfo.CodeInfo{}
	err = model.DB.Model(&c).Where("id = ?", companyAdminTestCtx.DeleteExpressId).First(&c).Error
	assert.Error(t, err)
	var cs []codeinfo.CodeInfo
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
 * 公司管理员获取价格集成测试
 */
func TestCompanyAdminRoleGetPrice(t *testing.T) {
	q := `
		 query getPriceQuery {
		  getPrice
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员获取集成测试")
	}
}

/**
 * 公司管理员编辑价格集成测试
 */
func TestCompanyAdminRoleEditPrice(t *testing.T) {
	q := `
		mutation editPriceMutation($price: Float!){
		  editPrice(price: $price) 
		}
	`
	price := 134.4578
	v := map[string]interface{}{
		"price": price,
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员获取集成测试")
	}
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
	c := configs.Configs{}
	model.DB.Model(&configs.Configs{}).
		Where("name = ? AND company_id = ?", configs.PRICE_NAME, me.CompanyId).
		First(&c)
	expectPrice, _ := strconv.ParseFloat(c.Value, 64)
	assert.Equal(t, expectPrice, price)
}

/**
 * 公司管理员获取设备列表集成测试
 */
func TestCompanyAdminGetDeviceList(t *testing.T) {
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
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
	// 断言响应的数据就是用户的公司名下的
	me, _ := GetUserByToken(companyAdminTestCtx.Token)
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
		// 用于编辑设备测试
		companyAdminTestCtx.EditDeviceId = int64(id)
	}
}

/**
 * 公司管理员编辑设备集成测试
 */
func TestCompanyAdminEditDevice(t *testing.T) {
	q := `
		mutation editDeviceMutation($input: EditDeviceInput!){
		  editDevice(input: $input) 
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"id":     companyAdminTestCtx.EditDeviceId,
			"isAble": false,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员编辑设备集成测试")
	}
	d := devices.Device{}
	err = d.GetDeviceSelfById(companyAdminTestCtx.EditDeviceId)
	assert.NoError(t, err)
	assert.False(t, d.IsAble)

	v = map[string]interface{}{
		"input": map[string]interface{}{
			"id":     companyAdminTestCtx.EditDeviceId,
			"isAble": true,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("failed:公司管理员编辑设备集成测试")
	}
	err = d.GetDeviceSelfById(companyAdminTestCtx.EditDeviceId)
	assert.NoError(t, err)
	assert.True(t, d.IsAble)
}

/**
 * 公司管理员获取型钢列表集成测试
 */
func TestCompanyAdminGetSteelList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}

/**
 * 公司管理员创建项目集成测试
 */
func TestCompanyAdminCreateProject(t *testing.T) {
	var totalLogs int64
	model.DB.Model(&logs.Logos{}).Count(&totalLogs)
	q := `
		query ($input: GetCompanyUserInput!) {
		  getCompanyUser(input: $input) {
			id
			name
		  }
		}
	`
	v := map[string]interface{}{
		"input": map[string]interface{}{
			"roleId": roles.RoleProjectAdminId,
		},
	}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("公司管理员创建项目集成测试")
	}
	users := res["getCompanyUser"].([]interface{})
	var userIds []int64
	for _, user := range users {
		u := user.(map[string]interface{})
		userIds = append(userIds, int64(u["id"].(float64)))
	}
	q = `
		mutation ($input: CreateProjectInput!) {
		  createProject(input: $input) {
			id
			city
			company{id name}
			endedAt
			leaderList { id name }
			city
			name
			remark
			startedAt
		  }
		}
	`
	s := time.Unix(time.Now().Unix()+1000, 0)
	timeStr := s.Format(time.RFC3339)
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"address":   "address_for_companyCreateProjectTest",
			"city":      "city_for_companyCreateProjectTest",
			"leaderIdS": userIds,
			"name":      "name_for_companyCreateProjectTest",
			"remark":    "remark_for_companyCreateProjectTest",
			"startAt":   timeStr,
		},
	}
	_, err = graphReqClient(q, v, roles.RoleCompanyAdmin)
	if err != nil {
		t.Fatal("公司管理员创建项目集成测试")
	}
	// 日志新增断言
	var currentTotalLogs int64
	model.DB.Model(&logs.Logos{}).Count(&currentTotalLogs)
	assert.Equal(t, totalLogs+1, currentTotalLogs)
}

/**
 * 公司管理员获取项目列表集成测试
 */
func TestCompanyAdminGetProjectList(t *testing.T) {
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}

/**
 * 公司管理员获取订单列表集成测试
 */
func TestCompanyAdminGetOrderList(t *testing.T) {
	q := `
			
# confirmOrder 确认订单类型
# ≈ 待确认订单类型
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
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}

/**
 * 公司管理员添加维修厂集成测试
 */
func TestCompanyAdminRoleCreateMaintenance(t *testing.T) {
	q := `
		mutation ($input: CreateMaintenanceInput!){
		   createMaintenance(input: $input) {
			id
			# 备注
			remark
			# 地址

			address
			# 维修厂名
			name
			# 管理员列表
			admin{
			  id
			   #名字
			  name
			  #电话
			  phone
			  # 微信 
			  wechat
			}
		  }
		}
	`
	v = map[string]interface{}{
		"input": map[string]interface{}{
			"address": "测试地址1",
			"name":    "测试名1",
			"uid":     5,
		},
	}
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}

/**
 * 公司管理员 编辑维修厂  集成测试
 */
func TestCompanyAdminRoleEditMaintenance(t *testing.T) {
	q := `
		mutation ($input: EditMaintenanceInput!) {
		  editMaintenance(input: $input) {
			 id
			# 备注
			remark
			# 地址

			address
			# 维修厂名
			name
			# 管理员列表
			admin{
			  id
			   #名字
			  name
			  #电话
			  phone
			  # 微信 
			  wechat
			}     
		  }
		}
	`
	_, err := graphReqClient(q, v, roles.RoleCompanyAdmin)
	assert.NoError(t, err)
}
