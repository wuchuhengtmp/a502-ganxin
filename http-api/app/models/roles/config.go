/**
 * @Desc    The roles is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/24
 * @Listen  MIT
 */
package roles

import (
	"fmt"
	"gorm.io/gorm"
	"io"
	"strconv"
)

type Roles struct {
	ID   int64  `json:"id"`
	Name string `json:"name" gorm:"comment:角色名"`
	Tag  Role   `json:"tag" gorm:"comment:角色标识"`
	gorm.Model
}

type Role string

const (
	//  超级管理员
	RoleAdmin Role = "admin"
	//  公司管理员
	RoleCompanyAdmin Role = "companyAdmin"
	//  仓库管理员
	RoleRepositoryAdmin Role = "repositoryAdmin"
	//  项目管理员
	RoleProjectAdmin Role = "projectAdmin"
	//  维修管理员
	RoleMaintenanceAdmin Role = "maintenanceAdmin"
)

func (e Role) IsValid() bool {
	switch e {
	case RoleAdmin, RoleCompanyAdmin, RoleRepositoryAdmin, RoleProjectAdmin, RoleMaintenanceAdmin:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
