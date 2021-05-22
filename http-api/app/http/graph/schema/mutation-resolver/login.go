/**
 * @Desc    The mutation_resolver is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/22
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"errors"
	"http-api/app/http/graph/model"
)

/**
 * 登录
 */
func (r *MutationResolver) Login (ctx context.Context, phone *string, password *string) (*model.LoginRes, error)   {
	err := errors.New("密码错误")
	return &model.LoginRes{

	}, err

}
