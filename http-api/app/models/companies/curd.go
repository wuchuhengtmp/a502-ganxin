/**
 * @Desc    The companies is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package companies

import "http-api/pkg/model"

func (b *Companies)Create () error  {
	db := model.DB
	return db.Model(b).Create(b).Error
}
