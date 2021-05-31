/**
 * @Desc    The companies is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package companies

import (
	"http-api/app/models/roles"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"time"
)

func (c *Companies)Create () error  {
	db := model.DB
	return db.Model(c).Create(c).Error
}

func GetAll() (cs []Companies) {
	db := model.DB
	db.Model(&Companies{}).Find(&cs)
	// 公司状态根据有郊期期限修正
	for i, company := range cs {
		if company.IsAble == true && company.EndedAt.Unix() >= time.Now().Unix() {
			company.IsAble = false
			db.Model(&company).Where("id = ?", company.ID).Updates(company)
			cs[i].IsAble = false
		}
	}

	return cs
}
/**
 * 获取公司的管理员
 */
func (c *Companies) GetAdmin() (user users.Users, err error) {
	db := model.DB
	err = db.Model(&user).
		Where("role_id = ? AND company_id = ?", roles.RoleCompanyAdminId, c.ID).
		First(&user).Error
	return user, err
}

/**
 *  有没有这家公司
 */
func (Companies) HasCompanyId(cid int64) (*Companies, error)  {
	 db := model.DB
	 cp := Companies{}
	 err := db.Model(&Companies{}).Where("id = ?", cid).First(&cp).Error
	 return &cp, err
}
