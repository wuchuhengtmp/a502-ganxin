/**
 * @Desc    The roles is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package roles

import (
	"fmt"
	"gorm.io/gorm"
	"http-api/pkg/model"
	"io"
	"strconv"
)

type Role struct {
	ID   int64       `json:"id"`
	Name string      `json:"name" gorm:"comment:角色名"`
	Tag  GraphqlRole `json:"tag" gorm:"comment:角色标识"`
	gorm.Model
}

func (Role)TableName() string  {
	return "roles"
}

func (r *Role) GetSelfById(id int64) error {
	db := model.DB
	return db.Model(r).Where("id = ?", id).First(r).Error
}

type GraphqlRole string

const (
	//  超级管理员
	RoleAdmin GraphqlRole = "admin"
	//  公司管理员
	RoleCompanyAdmin GraphqlRole = "companyAdmin"
	//  仓库管理员
	RoleRepositoryAdmin GraphqlRole = "repositoryAdmin"
	//  项目管理员
	RoleProjectAdmin GraphqlRole = "projectAdmin"
	//  维修管理员
	RoleMaintenanceAdmin GraphqlRole = "maintenanceAdmin"

	//  超级管理员Id
	RoleAdminId int64 = 1
	//  公司管理员Id
	RoleCompanyAdminId int64 = 2
	//  仓库管理员id
	RoleRepositoryAdminId int64 = 3
	//  项目管理员id
	RoleProjectAdminId int64 = 4
	//  维修管理员id
	RoleMaintenanceAdminId int64 = 5
)

// 角色标识映射角色id
var RoleTagMapId = map[string]int64{
	RoleAdmin.String():            RoleAdminId,
	RoleCompanyAdmin.String():     RoleCompanyAdminId,
	RoleRepositoryAdmin.String():  RoleRepositoryAdminId,
	RoleProjectAdmin.String():     RoleProjectAdminId,
	RoleMaintenanceAdmin.String(): RoleMaintenanceAdminId,
}

// 角色标识映射角色名
var RoleTagMapName = map[string]string{
	RoleAdmin.String():            "超级管理员",
	RoleCompanyAdmin.String():     "公司管理员",
	RoleRepositoryAdmin.String():  "仓库管理员",
	RoleProjectAdmin.String():     "项目管理员",
	RoleMaintenanceAdmin.String(): "维修管理员",
}


func (e GraphqlRole) IsValid() bool {
	switch e {
	case RoleAdmin, RoleCompanyAdmin, RoleRepositoryAdmin, RoleProjectAdmin, RoleMaintenanceAdmin:
		return true
	}
	return false
}

func (e GraphqlRole) String() string {
	return string(e)
}

func (e *GraphqlRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GraphqlRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e GraphqlRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

/**
 * 获取解决列表的解析器的响应数据列表
 */
func GetRolesGraphRes() ([]*Role, error) {
	db := model.DB
	var res []*Role
	var roles []Role
	err := db.
		Model(&Role{}).
		Where("tag In (?)", []GraphqlRole{RoleRepositoryAdmin, RoleProjectAdmin, RoleMaintenanceAdmin}).
		Find(&roles).Error
	if  err != nil {
		return res, err
	}
	for _, role := range roles {
		tmp := Role{}
		tmp.ID = role.ID
		tmp.Name = role.Name
		tmp.Tag =  role.Tag
		res = append(res, &tmp)
	}

	return res, nil
}
