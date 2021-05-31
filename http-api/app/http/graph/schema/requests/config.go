/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/28
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/files"
	"http-api/app/models/users"
	"regexp"
)

func init()  {
	// 定义验证手机规则
	govalidator.AddCustomRule("phone", func(field string, rule string, message string, value interface{}) error {
		const patter string = `^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199)\d{8}$`
		reg := regexp.MustCompile(patter)
		if reg.Match([]byte(value.(string))) {
			return nil
		} else {
			return fmt.Errorf("%s:%s 不是正确的手机号", field, value)
		}

		return nil
	})
	// 时间字串验证规则
	govalidator.AddCustomRule("time", func(field string, rule string, message string, value interface{}) error {
		_, err := helper.Str2Time(value.(string))
		if err != nil {
			return fmt.Errorf("%s:%s 不是正确的 类2006-01-02 15:04:05 时间格式", field, value)
		}

		return nil
	})

	// 是否存在这个文件验证规则
	govalidator.AddCustomRule("fileExist", func(field string, rule string, message string, value interface{}) error {
		id := int64(value.(int))
		file := files.File{
			ID: id,
		}
		if !file.IsExist() {
			return fmt.Errorf("%s:%d 该文件不存在", field, id)
		}

		return nil
	})

	// 用户手机号不能存在
	govalidator.AddCustomRule("not_user_phone_exists", func(field string, rule string, message string, value interface{}) error {
		userModel := users.Users{}
		if userModel.IsPhoneExists(value.(string)) {
			return fmt.Errorf("%s:%d 手机号已存在", field)
		}

		return nil
	})
}
