/**
 * @Desc    The users is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/28
 * @Listen  MIT
 */
package users

import "http-api/pkg/model"

func (u *Users)Create () error  {
	db := model.DB
	return db.Model(u).Create(u).Error
}
