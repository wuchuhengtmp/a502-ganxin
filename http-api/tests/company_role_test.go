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
	"github.com/magiconair/properties/assert"
	"http-api/app/models/codeinfo"
	"http-api/app/models/roles"
	"http-api/pkg/model"
	"http-api/seeders"
	"math/rand"
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
			adminName
			total
			adminWechat
			adminPhone
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
	hasError(t, err)
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
	hasError(t, err)
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
			"isDefault": isDefault,
		},
	}
	if _, err := graphReqClient(q, v, roles.RoleCompanyAdmin); err != nil {
		t.Fatal(err.Error())
	}
	newCodeInfo := codeinfo.CodeInfo{}
	err := model.DB.Model(&codeinfo.CodeInfo{}).
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
	v := map[string]interface{} {}
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin);
	if  err != nil {
		t.Fatal("failed:公司管理员获取制造商列表集成测试")
	}
	assertCompanyIdForGetManufacturers(t, res, companyAdminTestCtx.Token)
}


