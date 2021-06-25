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
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/util/helper"
	"http-api/app/models/codeinfo"
	"http-api/app/models/devices"
	"http-api/app/models/files"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/repositories"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
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
	// 是否是公司名下的型钢
	govalidator.AddCustomRule("isCompanySteel", func(field string, rule string, message string, value interface{}) error {
		me, err := getUserByRule(rule)
		if err != nil { return err }
		r := steels.Steels{}
		err = model.DB.
			Model(&r).
			Where("company_id = ? AND id = ?", me.CompanyId, value).First(&r).
			Error
		if err != nil {
			return fmt.Errorf("公司中没有这个型钢")
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

// 验证可待出库的订单型钢详情验证器的验证步骤合集
type ValidateGetProject2WorkshopDetailRequestSteps map[string]int64

// 有没有这个订单
func (ValidateGetProject2WorkshopDetailRequestSteps) checkHasOrder(ctx context.Context, orderId int64) error {
	me := auth.GetUser(ctx)
	o := orders.Order{}
	err := model.DB.Model(&o).Where("id = ?", orderId).
		Where("company_id = ?", me.CompanyId).
		First(&o).
		Error
	if err != nil {
		return fmt.Errorf("没有这个订单id: %d", orderId)
	}

	return nil
}

// 检验订单状态 只能是确认或部分发货才行
func (ValidateGetProject2WorkshopDetailRequestSteps) checkOrderState(ctx context.Context, orderId int64) error {
	o := orders.Order{}
	if e := model.DB.Model(&o).Where("id = ?",orderId).First(&o).Error; e != nil {
		return fmt.Errorf("没有这个订单id:%d", orderId)
	}
	if o.State != orders.StateConfirmed && o.State != orders.StatePartOfReceipted {
		return fmt.Errorf("当前订单状态为:%s, 不能接着发货", orders.StateMapDesc[o.State])
	}

	return nil
}

// 获取订单各规格的需求量上限
func (ValidateGetProject2WorkshopDetailRequestSteps) GetOrderSpecificationGroupTotal(ctx context.Context, orderId int64) (map[string]int64, error) {
	list := make(map[string]int64)
	var osList []*order_specification.OrderSpecification
	if err := model.DB.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", orderId).Find(&osList).Error; err != nil {
		return list, nil
	}
	for _, item := range osList {
		oss := order_specification_steel.OrderSpecificationSteel{}
		var existsTotal int64
		if err := model.DB.Model(&oss).Where("order_specification_id = ?", item.Id).Count(&existsTotal).Error; err != nil {
			return list, err
		}
		list[item.Specification] = item.Total - existsTotal
	}

	return list, nil
}

/**
 * 检验是否有冗余识别码
 */
func (ValidateGetProject2WorkshopDetailRequestSteps) CheckRedundancyIdentification(list []string) error {
	identificationMapTotal := make(map[string]int64)
	for _, item := range list {
		if _, ok := identificationMapTotal[item]; ok {
			return fmt.Errorf("识别码出现冗余，%s 不能输入多个同样的", item)
		} else {
			identificationMapTotal[item] = 1
		}
	}

	return nil
}

/*
 * 识别码不能为空
 */
func (ValidateGetProject2WorkshopDetailRequestSteps) CheckIdentificationListMustBeEmpty(ctx context.Context, identifierList []string) error {
	if len(identifierList) == 0 {
		return fmt.Errorf("识别码列表不能为空")
	}

	return nil
}

func (ValidateGetProject2WorkshopDetailRequestSteps) CheckSteelList(ctx context.Context, orderId int64, identifierList []string) error {
	me := auth.GetUser(ctx)
	// 订单规格合集
	var orderSpecificationList []*order_specification.OrderSpecification
	orderSpecificationSpecificationMapTotal := make(map[string]int64) // 当前同一规格统计量 用于比较上限
	var orderSpecificationIdList []int64 // 订单要求的规格id集合，用于检验型钢的规格是否在这个合集中
	err := model.DB.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", orderId).
		Find(&orderSpecificationList).
		Error
	if err != nil {
		return err
	}
	for _, item := range orderSpecificationList {
		orderSpecificationIdList = append(orderSpecificationIdList, item.SpecificationId)
	}
	// 获取订单各规格需求上限
	osl, err := ValidateGetProject2WorkshopDetailRequestSteps{}.GetOrderSpecificationGroupTotal(ctx, orderId)
	if err != nil {
		return nil
	}
	// 检验每根型钢
	for _, identification := range identifierList {
		s := steels.Steels{}
		// 检验型钢状态能否满足订单要求
		err := model.DB.Model(&steels.Steels{}).
			Where("identifier = ?", identification).
			//Where("state = ?", steels.StateInStore).
			Where("company_id = ?", me.CompanyId).
			First(&s).
			Error
		if err != nil {
			return fmt.Errorf("仓库中没有 %s 标识码的型钢在仓库中", identification)
		}
		// 检验型钢状态
		if s.State != steels.StateInStore {
			return fmt.Errorf("识别码为%s的型钢当前状态为:%s, 不能出库", identification, steels.StateCodeMapDes[s.State])
		}
		// 检验型钢的规格能否满足订单的要求
		if err := func() error{
			for _, specificationId := range orderSpecificationIdList {
				if specificationId == s.SpecificationId  {
					return nil
				}
			}
			return fmt.Errorf("订单中,要求的规格id为:%v, 而标识码的为%s的型钢的规格id为%d, 并不能满足订单的要求", orderSpecificationIdList, identification, s.SpecificationId)
		}(); err != nil {
			return err
		}
		specificationInstance, err := s.GetSpecification()
		if err != nil {
			return fmt.Errorf("型钢规格不存在 id:%d ，请联系管理员", identification)
		}
		// 检验当前规格型钢的数量是否超过订单要求的上限
		key := specificationInstance.GetSelfSpecification()
		orderSpecificationSpecificationMapTotal[key] += 1
		// 上限比较
		if orderSpecificationSpecificationMapTotal[key] > osl[key] {
			return fmt.Errorf("当前规格%s， 已经超过订单要求的%d 数量了", key, osl[key])
		}
	}

	return nil
}
/**
 * 检验规格
 */
func (ValidateGetProject2WorkshopDetailRequestSteps) CheckSpecification(ctx context.Context, orderId int64, specificationId *int64) error {
	if specificationId!= nil {
		err := model.DB.
			Model(&order_specification.OrderSpecification{}).
			Where("order_id = ?", orderId).
			Where("specification_id = ?", *specificationId).
			First(&order_specification.OrderSpecification{}).
			Error
		if err != nil {
			return fmt.Errorf("订单中没有id为: %d 的规格", *specificationId)
		}
	}

	return nil
}
// 获取项目规格列表验证器验证步骤
type ValidateGetProjectSpecificationDetailRequestSteps struct {}

/**
 * 项目的管理员是不是我
 */
func (v *ValidateGetProjectSpecificationDetailRequestSteps)CheckProjectLeader(ctx context.Context, projectId int64) error  {
	me := auth.GetUser(ctx)
	projectTable := projects.Projects{}.TableName()
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	projectItem := projects.Projects{}
	err := model.DB.Model(&projectItem).
		Select(fmt.Sprintf("%s.*", projectTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", projectLeaderTable, projectLeaderTable, projectTable)).
		Where(fmt.Sprintf("%s.uid = ?", projectLeaderTable), me.Id).
		Where(fmt.Sprintf("%s.id = ?", projectTable), projectId).
		First(&projectItem).
		Error
	if err != nil {
		return fmt.Errorf("项目%d, 管理员不是您，您当前无权操作", projectId)
	}

	return nil
}

/**
 * 项目是否存在
 */
func (*ValidateGetProjectSpecificationDetailRequestSteps)CheckProjectExists(ctx context.Context, projectId int64) error  {
	projectItem := projects.Projects{}
	me := auth.GetUser(ctx)
	err := model.DB.Model(&projectItem).Where("id = ? AND company_id = ?", projectId, me.CompanyId).First(&projectItem).Error
	if err != nil {
		return fmt.Errorf("没有项目id为：%d 的项目", projectId)
	}

	return nil
}

