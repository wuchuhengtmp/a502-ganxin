/**
 * @Desc    订单模型
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/26
 * @Listen  MIT
 */
package orders

import (
	"fmt"
	"gorm.io/gorm"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/specificationinfo"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"time"
)

type Order struct {
	Id               int64     `json:"id"`
	ProjectId        int64     `json:"ProjectId" gorm:"comment:项目id"`
	RepositoryId     int64     `json:"repositoryId" gorm:"comment: 仓库id"`
	State            StateCode `json:"state" gorm:"comment:订单状态 待确认200 已确认300 已拒绝400 已发货500 待收货600 已收货(部分)700 已收货全部800 已归库900"`
	PartList         string    `json:"partList" gorm:"comment:配件清单"`
	CreateUid        int64     `json:"createUid" gorm:"comment:创建人"`
	ConfirmedUid     int64     `json:"confirmedUid" gorm:"comment:确认人"`
	OrderNo          string    `json:"orderNo" gorm:"comment:订单编号"`
	Remark           string    `json:"remark" gorm:"comment:备注"`
	CompanyId        int64     `json:"companyId" gorm:"comment:公司id"`
	ExpectedReturnAt time.Time `json:"expectedReturnAt" gorm:"comment:预计归还时间"`
	ReceiveAt        time.Time `json:"receiveAt" gorm:"comment:收货时间"`
	ConfirmedAt      time.Time `json:"confirmAt" gorm:"comment:确认时间"`
	RejectReason     string    `json:"rejectReason" gorm:"comment:描绘原因"`
	gorm.Model
}

func (Order) TableName() string {
	return "orders"
}

type StateCode = int64

const (
	StateToBeConfirmed   StateCode = 200 // 待确认
	StateConfirmed       StateCode = 300 // 已确认
	StateRejected        StateCode = 400 // 已拒绝
	StateSend            StateCode = 500 // 已发货
	StatePartOfReceipted StateCode = 700 // 部分收货
	StateReceipted       StateCode = 800 // 收货
)

/**
 *  状态码映射说明
 */
var StateMapDesc = map[StateCode]string{
	StateToBeConfirmed:   "待确认",
	StateConfirmed:       "已确认",
	StateRejected:        "已拒绝",
	StateSend:            "已发货",
	StatePartOfReceipted: "部分收货",
	StateReceipted:       "收货",
}

/**
 * 根据物流公司获取订单列表
 */
func (*Order) GetOrdersByExpressId(expressId int64) (os []*Order, err error) {
	err = model.DB.Model(&Order{}).Where("express_company_id = ?", expressId).Find(&os).Error

	return
}

/**
 *  获取确认订单型钢的数量
 */
func GetConfirmSteelTotalBySpecificationId(specificationId int64) (int64, error) {
	o := Order{}
	oss := order_specification_steel.OrderSpecificationSteel{}
	os := order_specification.OrderSpecification{}
	st := steels.Steels{}
	var confirmTotal int64
	err := model.DB.Model(&oss).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", os.TableName(), os.TableName(), oss.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", o.TableName(), o.TableName(), os.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", st.TableName(), st.TableName(), oss.TableName())).
		Where(fmt.Sprintf("%s.specification_id = %d", st.TableName(), specificationId)).
		Where(fmt.Sprintf("%s.state = %d", o.TableName(), StateConfirmed)).
		Count(&confirmTotal).Error

	return confirmTotal, err
}

/**
* 获取确认订单的重量
 */
func GetConfirmSteelTotalWeightBySpecificationId(specificationId int64) (float64, error) {
	o := Order{}
	oss := order_specification_steel.OrderSpecificationSteel{}
	os := order_specification.OrderSpecification{}
	st := steels.Steels{}
	var osList []*order_specification.OrderSpecification
	err := model.DB.Model(&oss).
		Select(fmt.Sprintf("%s.*", os.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", os.TableName(), os.TableName(), oss.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", o.TableName(), o.TableName(), os.TableName())).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.steel_id", st.TableName(), st.TableName(), oss.TableName())).
		Where(fmt.Sprintf("%s.specification_id = %d", st.TableName(), specificationId)).
		Where(fmt.Sprintf("%s.state = %d", o.TableName(), StateConfirmed)).
		Find(&osList).Error
	var totalWeight float64
	for _, orderSpecificationItem := range osList {
		spe := specificationinfo.SpecificationInfo{}
		err := model.DB.Model(&spe).
			Where("id = ?", orderSpecificationItem.SpecificationId).
			First(&spe).Error
		if err != nil {
			return totalWeight, err
		}
		totalWeight += spe.Weight * float64(orderSpecificationItem.Total)
	}

	return totalWeight, err
}

func (o *Order) GetSelf() (err error) {
	err = model.DB.Model(&Order{}).Where("id = ?", o.Id).First(&o).Error

	return
}

/**
 * 获取订单上型钢的数量
 */
func GetTotal(tx *gorm.DB, o *Order) (int64, error) {
	var totalSteels struct {
		TotalSteels int64
	}
	err := tx.Model(&order_specification.OrderSpecification{}).
		Select("sum(total) as TotalSteels").
		Where("order_id = ?", o.Id).
		Scan(&totalSteels).Error

	return totalSteels.TotalSteels, err
}

/**
 * 获取订单重量
 */
func GetWeight(tx *gorm.DB, o *Order) (float64, error) {
	var totalWeight float64
	var oss []*order_specification.OrderSpecification
	tx.Model(&order_specification.OrderSpecification{}).Where("order_id = ?", o.Id).Find(&oss)
	for _, specificationItem := range oss {
		spec := specificationinfo.SpecificationInfo{}
		if err := tx.Model(&spec).Where("id = ?", specificationItem.SpecificationId).First(&spec).Error; err != nil {
			return 0, err
		}
		totalWeight += spec.Weight * float64(specificationItem.Total)
	}

	return totalWeight, nil
}

//""" 获取订单详情(用于管理后台)参数 """
type GetOrderDetailForBackEndRes struct {
	//""" 订单列表 """
	List []*order_specification.OrderSpecification
	//""" 数量 """
	Total int64
	//""" 重量 """
	Weight float64
	// 型钢数量
	SteelTotal int64
}
