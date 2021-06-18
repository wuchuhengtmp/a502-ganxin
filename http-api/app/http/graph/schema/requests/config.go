/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
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
	"http-api/app/models/codeinfo"
	"http-api/app/models/devices"
	"http-api/app/models/files"
	"http-api/app/models/orders"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"regexp"
	"strconv"
	"strings"
)

func init() {
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
		id := value.(int64)
		file := files.File{
			ID: id,
		}
		if !file.IsExist() {
			return fmt.Errorf("%s:%d 该文件不存在", field, id)
		}

		return nil
	})
	// 是否存在这个用户验证规则
	govalidator.AddCustomRule("userExist", func(field string, rule string, message string, value interface{}) error {
		uid := value.(int64)
		user := users.Users{}
		err := user.GetSelfById(uid)
		if err != nil {
			return fmt.Errorf("%s:%d 该用户不存在", field, uid)
		}

		return nil
	})
	// 用户手机号不能存在
	govalidator.AddCustomRule("not_user_phone_exists", func(field string, rule string, message string, value interface{}) error {
		userModel := users.Users{}
		if userModel.IsPhoneExists(value.(string)) {
			return fmt.Errorf("%s:%d 手机号已存在", field, value)
		}

		return nil
	})
	// 是否大于0
	govalidator.AddCustomRule("isGreaterZero", func(field string, rule string, message string, value interface{}) error {
		v := value.(float64)
		if v <= 0 {
			return fmt.Errorf("%s:%f 必须大于0", field, v)
		}

		return nil
	})

	// 规格表的id是否存在
	govalidator.AddCustomRule("isSpecificationId", func(field string, rule string, message string, value interface{}) error {
		v := value.(int64)
		s := specificationinfo.SpecificationInfo{ID: v}
		if err := s.GetSelf(); err != nil {
			return fmt.Errorf("%s:%d 没有这个规格记录", field, v)
		}

		return nil
	})

	// 码表的id是否存在
	govalidator.AddCustomRule("isCodeInfoId", func(field string, rule string, message string, value interface{}) error {
		v := value.(int64)
		c := codeinfo.CodeInfo{ID: v}
		if err := c.GetSelf(); err != nil {
			return fmt.Errorf("%s:%d 没有这个码表记录", field, v)
		}

		return nil
	})

	// 设备的id是否存在
	govalidator.AddCustomRule("isDeviceId", func(field string, rule string, message string, value interface{}) error {
		v := value.(int64)
		c := devices.Device{ID: v}
		if err := c.GetDeviceSelfById(v); err != nil {
			return fmt.Errorf("%s:%d 没有这个设备记录", field, v)
		}

		return nil
	})
	// 长度
	govalidator.AddCustomRule("minLen", func(field string, rule string, message string, value interface{}) error {
		v := value.(int64)
		fmt.Println(v)

		return nil
	})
	// 是否是公司仓库
	govalidator.AddCustomRule("isCompanyRepository", func(field string, rule string, message string, value interface{}) error {
		me, err := getUserByRule(rule)
		if err == nil { return err }
		r := repositories.Repositories{}
		err = model.DB.
			Model(&repositories.Repositories{}).
			Where("company_id = ? AND id = ?", me.CompanyId, value).First(&r).
			Error
		if err != nil {
			return fmt.Errorf("公司中没有这个仓库")
		}

		return nil
	})

	// 是否是公司名下的规格
	govalidator.AddCustomRule("isCompanySpecification", func(field string, rule string, message string, value interface{}) error {
		me, err := getUserByRule(rule)
		if err != nil { return err }
		r := repositories.Repositories{}
		err = model.DB.
			Model(&specificationinfo.SpecificationInfo{}).
			Where("company_id = ? AND id = ?", me.CompanyId, value).First(&r).
			Error
		if err != nil {
			return fmt.Errorf("公司中没有这个规格")
		}

		return nil
	})
	// 是否是公司名下的项目
	govalidator.AddCustomRule("isCompanyProject", func(field string, rule string, message string, value interface{}) error {
		me, err := getUserByRule(rule)
		if err != nil { return err }
		r := projects.Projects{}
		err = model.DB.
			Model(&r).
			Where("company_id = ? AND id = ?", me.CompanyId, value).First(&r).
			Error
		if err != nil {
			return fmt.Errorf("公司中没有这个项目")
		}

		return nil
	})
	// 是否是公司下的订单
	govalidator.AddCustomRule("isCompanyOrder", func(field string, rule string, message string, value interface{}) error {
		me, err := getUserByRule(rule)
		if err != nil { return err }
		r := orders.Order{}
		err = model.DB.
			Model(&r).
			Where("company_id = ? AND id = ?", me.CompanyId, value).First(&r).
			Error
		if err != nil {
			return fmt.Errorf("公司中没有这个订单")
		}

		return nil
	})
}

func getUserByRule (rule string) (u *users.Users,  err error ) {
	uid, _ := strconv.ParseInt( strings.SplitAfter(rule, ":")[1], 10, 16)
	me := users.Users{}
	if err := me.GetSelfById(uid); err != nil { return nil, fmt.Errorf("没有这个用户") }

	return &me, nil
}

/**
 *  验证
 */
func Validate(opts govalidator.Options) error {
	errs := govalidator.New(opts).ValidateStruct()
	if len(errs) > 0 {
		for _, fieldErrors := range errs {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}
