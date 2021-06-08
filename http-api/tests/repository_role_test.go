/**
 * @Desc    仓库管理员角色集成测试
 * @Author  wuchuheng<wuchuheng@163.com>
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
	"http-api/app/models/roles"
	"http-api/pkg/model"
	"http-api/seeders"
	"math/rand"
	"testing"
	"time"
)

// 仓库管理员测试上下文
var repositoryAdminTestCtx = struct {
	Token    string
	Username string
	Password string
	// 用于个性规格记录
	EditSpecificationId   int64
	DeleteSpecificationId int64
	// 用于编辑的材料商家id
	EditMaterialId int64
	// 用于删除材料商家
	DeleteMaterialId int64
	// 用于编辑制造商家
	EditManufacturerId int64
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
	res, err := graphReqClient(q, v, roles.RoleCompanyAdmin);
	if  err != nil {
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
		assert.Len(t,  cs, 1)
	}
	// 保存id用于编辑制造商测试
	data := res["createManufacturer"].(map[string]interface{})
	id := data["id"].(float64)
	repositoryAdminTestCtx.EditManufacturerId = int64(id)
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
	v := map[string]interface{} {}
	res, err := graphReqClient(q, v, roles.RoleRepositoryAdmin);
	if  err != nil {
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
			"id": repositoryAdminTestCtx.EditManufacturerId,
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
		ID:repositoryAdminTestCtx.EditManufacturerId,
	}
	_ = c.GetSelf()
	assert.Equal(t, name, c.Name)
	assert.Equal(t, isDefault, c.IsDefault)
	assert.Equal(t, remark, c.Remark)
	// 只能有一个默认选项
	var cs  []*codeinfo.CodeInfo
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
