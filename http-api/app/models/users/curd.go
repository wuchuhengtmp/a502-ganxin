/**
 * @Desc    The users is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/28
 * @Listen  MIT
 */
package users

import "http-api/pkg/model"


type WechatModel struct {
	WcOpenId     string `json:"wc_openid" gorm:"comment:微信openid"`      // 用户唯一标识
	WcSessionKey string `json:"wc_session_key" gorm:"comment:微信会话i"` // 会话密钥
	WcUnionId    string `json:"wc_unionid" gorm:"comment:微信unionI"`     // 用户在开放平台的唯一标识符，在满足UnionID下发条件的情况下会返回
}
/**
 * 有没有这个用户
 */
func (WechatModel) IsUserByOpenId(wcOpenId string) bool {
	var user = Users{}
	err := model.DB.Model(&Users{}).Where("wc_open_id", wcOpenId).First(&user).Error
	return err == nil
}

func (w *WechatModel) AddUser() (Users, error) {
	user := Users{}
	err := model.DB.Model(&Users{}).Create(&user).Error
	return user, err
}

func (WechatModel) GetUserByOpenId(openId string) (user Users, err error)  {
	err = model.DB.Model(&user).Where("wc_open_id", openId).First(&user).Error
	return user, err
}
