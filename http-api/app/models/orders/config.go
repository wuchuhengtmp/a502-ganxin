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
	State            int       `json:"state" gorm:"订单状态 待确认200 已确认300 已拒绝400 已发货500 待收货600 已收货(部分)700 已收货全部800 已归库"`
	ExpectedReturnAt time.Time `json:"expectedReturnAt" gorm:"comment:预计归还时间"`
	PartList         string    `json:"partList" gorm:"comment:配件清单"`
	CreateUid        int64     `json:"createUid" gorm:"comment:创建人"`
	ConfirmedAt      time.Time `json:"confirmAt" gorm:"comment:确认时间"`
	ReceiveUid       int64     `json:"receiveUid" gorm:"comment:收货人id"`
	ReceiveAt        time.Time `json:"receiveAt" gorm:"comment:收货时间"`
	ExpressCompanyId int64     `json:"expressCompanyId" gorm:"comment:快递公司(码表id)"`
	ExpressNo        string    `json:"expressNo" gorm:"comment:物流号"`
	OrderNo          string    `json:"orderNo" gorm:"comment:订单编号"`
	Remark           string    `json:"remark" gorm:"comment:备注"`
	gorm.Model
}

func (Order) TableName() string {
	return "orders"
}

const (
	StateToBeConfirmed   = 200 // 待确认
	StateConfirmed       = 300 // 已确认
	StateRejected        = 400 // 已拒绝
	StateShipped         = 500 // 已发货
	StatePartOfReceipted = 700 // 部分收货
	StateReceipted       = 800 // 收货
)

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