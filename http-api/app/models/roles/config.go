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

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name" gorm:"comment:角色名"`
	Tag  GraphqlRole   `json:"tag" gorm:"comment:角色标识"`
	gorm.Model
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
)

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
