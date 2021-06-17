/**
 * @Desc    无角色鉴权限制的公共接口集成测试
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package tests

import (
	"testing"
)

func TestGetRoleList(t *testing.T) {
	q := `
		query getRoleListQuery {
		  getRoleList{
			id
			name
			tag
		  }
		}
	`
	v := map[string]interface{}{}
	_, err := graphReqClient(q, v, "")
	hasError(t, err)
}