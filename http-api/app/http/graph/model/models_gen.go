// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"http-api/app/models/roles"
	"io"
	"strconv"
	"time"
)

//  确认订单或打拒绝
type ConfirmOrderInput struct {
	//  订单id
	ID int64 `json:"id"`
	//  是否通过
	IsAccess bool `json:"isAccess"`
}

//  创建公司参数
type CreateCompanyInput struct {
	//  公司名
	Name string `json:"name"`
	//  公司名称拼写简写
	PinYin string `json:"pinYin"`
	//   宗旨
	Symbol string `json:"symbol"`
	//  logo 文件Id
	LogoFileID int64 `json:"logoFileId"`
	//  App 背景图片Id
	BackgroundFileID int64 `json:"backgroundFileId"`
	//  账号状态
	IsAble bool `json:"isAble"`
	//  公司的电话
	Phone string `json:"phone"`
	//  公司的微信
	Wechat string `json:"wechat"`
	//  开始时间
	StartedAt string `json:"startedAt"`
	//  结束时间
	EndedAt string `json:"endedAt"`
	//  管理员名称
	AdminName string `json:"adminName"`
	//  管理员手机
	AdminPhone string `json:"adminPhone"`
	//  管理员密码
	AdminPassword string `json:"adminPassword"`
	//  管理员微信
	AdminWechat string `json:"adminWechat"`
	//  管理员头像Id
	AdminAvatarFileID int64 `json:"adminAvatarFileId"`
}

//  添加用户信息需要的信息
type CreateCompanyUserInput struct {
	Name     string              `json:"name"`
	Phone    string              `json:"phone"`
	Role     CreateInputUserRole `json:"role"`
	Wechat   string              `json:"wechat"`
	AvatarID int64               `json:"avatarId"`
	Password string              `json:"password"`
}

//  添加物流商需要的参数
type CreateExpressInput struct {
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//  创建维修厂参数
type CreateMaintenanceInput struct {
	//  维修员id
	UID int64 `json:"uid"`
	//  维修厂名称
	Name string `json:"name"`
	//  地址
	Address string `json:"address"`
	//  备注
	Remark *string `json:"remark"`
}

//  添加制造商参数
type CreateManufacturerInput struct {
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//  添加材料商参数
type CreateMaterialManufacturerInput struct {
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//   创建需求单
type CreateOrderInput struct {
	//  项目ID
	ProjectID int64 `json:"projectId"`
	//  出货仓库ID
	RepositoryID int64 `json:"repositoryId"`
	//  预计归还时间
	ExpectedReturnAt time.Time `json:"expectedReturnAt"`
	//  备注
	Remark string `json:"remark"`
	//  配件清单
	PartList string `json:"partList"`
	//  型钢列表
	SteelList []*CreateOrderSteelInput `json:"steelList"`
}

//  创建需求单的指定型钢单项参数
type CreateOrderSteelInput struct {
	//  数量
	Total int64 `json:"total"`
	//  规格ID
	SpecificationID int64 `json:"specificationId"`
}

//  创建项目需要的参数
type CreateProjectInput struct {
	//  城市
	City string `json:"city"`
	//  项目名
	Name string `json:"name"`
	//  地址
	Address string `json:"address"`
	//  多个负责人ids
	LeaderIDS []int64 `json:"leaderIdS"`
	//  备注
	Remark string `json:"remark"`
	//  开始时间
	StartAt time.Time `json:"startAt"`
}

//   创建仓库需要提交的参数
type CreateRepositoryInput struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	RepositoryAdminID int64  `json:"repositoryAdminId"`
	Remark            string `json:"remark"`
	PinYin            string `json:"pinYin"`
}

//  创建规格需要提交的参数
type CreateSpecificationInput struct {
	Type      string  `json:"type"`
	Length    float64 `json:"length"`
	Weight    float64 `json:"weight"`
	IsDefault bool    `json:"isDefault"`
}

//  型钢入库需要的参数
type CreateSteelInput struct {
	//  识别码
	IdentifierList []string `json:"identifierList"`
	//  当前存放的仓库id
	RepositoryID int64 `json:"repositoryId"`
	//  规格表id
	SpecificationID int64 `json:"specificationId"`
	//  料商类型id
	MaterialManufacturerID int64 `json:"materialManufacturerId"`
	//  制造商(生产商)id
	ManufacturerID int64 `json:"manufacturerId"`
	//  生产时间
	ProducedDate time.Time `json:"producedDate"`
}

type DelMaintenanceInput struct {
	//  工厂id
	ID int64 `json:"id"`
}

//  修改公司参数
type EditCompanyInput struct {
	//  公司ID
	ID int64 `json:"id"`
	//  公司名
	Name string `json:"name"`
	//  公司名称拼写简写
	PinYin string `json:"pinYin"`
	//   宗旨
	Symbol string `json:"symbol"`
	//  logo 文件Id
	LogoFileID int64 `json:"logoFileId"`
	//  App 背景图片Id
	BackgroundFileID int64 `json:"backgroundFileId"`
	//  账号状态
	IsAble bool `json:"isAble"`
	//  公司的电话
	Phone string `json:"phone"`
	//  公司的微信
	Wechat string `json:"wechat"`
	//  开始时间
	StartedAt string `json:"startedAt"`
	//  结束时间
	EndedAt string `json:"endedAt"`
	//  管理员名称
	AdminName string `json:"adminName"`
	//  管理员手机
	AdminPhone string `json:"adminPhone"`
	//  管理员密码
	AdminPassword *string `json:"adminPassword"`
	//  管理员微信
	AdminWechat string `json:"adminWechat"`
	//  管理员头像Id
	AdminAvatarFileID int64 `json:"adminAvatarFileId"`
}

//  编辑公司人员需要的数据
type EditCompanyUserInput struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	RoleID int64  `json:"roleId"`
	// " 是否启用
	IsAble bool `json:"isAble"`
}

type EditDeviceInput struct {
	ID     int64 `json:"id"`
	IsAble bool  `json:"isAble"`
}

//  编辑物流商需要的参数
type EditExpressInput struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//  编辑维修厂参数
type EditMaintenanceInput struct {
	ID int64 `json:"id"`
	//  维修厂名
	Name string `json:"name"`
	//  地址
	Address string `json:"address"`
	//  管理员id
	AdminIDList []int64 `json:"adminIdList"`
	//  备注
	Remark *string `json:"remark"`
	//  是否启用
	IsAble *bool `json:"isAble"`
}

//  编辑制作商需要的参数
type EditManufacturerInput struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//  编辑材料商需要的参数
type EditMaterialManufacturerInput struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Remark    string `json:"remark"`
	IsDefault bool   `json:"isDefault"`
}

//  修改规格需要提交的参数
type EditSpecificationInput struct {
	ID        int64   `json:"id"`
	Type      string  `json:"type"`
	Length    float64 `json:"length"`
	Weight    float64 `json:"weight"`
	IsDefault bool    `json:"isDefault"`
}

//  型钢入厂输入
type EnterMaintenanceInput struct {
	//  识别码
	Identifier string `json:"identifier"`
}

type ErrCodes struct {
	//  错误码编号
	Code int64 `json:"code"`
	//  错误码用途说明
	Desc string `json:"desc"`
}

type FileItem struct {
	//  文件ID
	ID int64 `json:"id"`
	//  文访问链接
	URL string `json:"url"`
}

//  获取待修改的仓库型钢详情参数
type Get2BeChangedRepositorySteelDetailInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
}

//  获取用于修改的仓库型钢参数
type Get2BeChangedRepositorySteelInput struct {
	//  识别码
	Identifier string `json:"identifier"`
}

//  获取型钢维修出库详情参数
type Get2BeMaintainSteelDetailInput struct {
	//  识别码
	IdentifierList []string `json:"identifierList"`
	//  规格 ID
	SpecificationID *int64 `json:"specificationId"`
}

//  获取型钢维修出库参数
type Get2BeMaintainSteelInput struct {
	//  识别码
	Identifier string `json:"identifier"`
}

//  获取用于报废的仓库型钢详情参数
type Get2BeScrapRepositorySteelDetailInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
}

//  获取用于报废的仓库型钢参数
type Get2BeScrapRepositorySteelInput struct {
	//  型钢标识码
	Identifier string `json:"identifier"`
}

type GetCompanyUserInput struct {
	//  角色id
	RoleID *int64 `json:"roleId"`
}

//  项目归库的型钢查询参数
type GetEnterRepositorySteelDetailInput struct {
	//  识别码
	Identifier string `json:"identifier"`
	//  项目列表
	ProjectID int64 `json:"projectId"`
}

//  获取项目最大安装码需要的参数
type GetMaxLocationCodeInput struct {
	//  型钢识别码
	Identifier string `json:"identifier"`
}

//  获取多个型钢详情的请求参数
type GetMultipleSteelDetailInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
}

//  获取一个型钢详情请求参数
type GetOneSteelDetailInput struct {
	Identifier string `json:"identifier"`
}

type GetOrderListInput struct {
	QueryType *GetOrderListInputType `json:"queryType"`
}

//  获取订单型钢详情需要的参数
type GetOrderSteelDetailInput struct {
	//  识别码
	Identifier string `json:"identifier"`
}

//  获取型钢单根型钢出场详情参数
type GetOutOfWorkshopProjectSteelDetail struct {
	//  标识码
	Identifier string `json:"identifier"`
	//  项目id
	ProjectID int64 `json:"projectId"`
}

//  获取要入场的订单需要的请求参数
type GetProjectOrder2WorkshopDetailInput struct {
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
	//  订单id
	OrderID int64 `json:"orderId"`
}

//  获取项目规格列表需要的请求参数
type GetProjectSpecificationDetailInput struct {
	//  项目id
	ProjectID int64 `json:"projectId"`
}

//  待修改武钢信息请求参数
type GetProjectSteel2BeChangeInput struct {
	//  识别码
	Identifier string `json:"identifier"`
}

//  获取项目型钢详情需要的参数
type GetProjectSteelDetailInput struct {
	//  项目id
	ProjectID int64 `json:"projectId"`
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
	//  型钢状态
	State *int64 `json:"state"`
}

type GetRepositoryOverviewInput struct {
	//  仓库id
	ID int64 `json:"id"`
	//  规格ID
	SpecificationID *int64 `json:"specificationId"`
}

//  获取仓库型钢信息参数
type GetRepositorySteelInput struct {
	//  仓库id
	ReposirotyID int64 `json:"reposirotyId"`
	//  规格尺寸
	SpecificationID *int64 `json:"specificationId"`
	//  当前状态
	State *int64 `json:"state"`
}

//  获取待归库详情参数
type GetToBeEnterRepositoryDetailInput struct {
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
	//  型钢状态
	State *int64 `json:"state"`
	//  项目ID
	ProjectID int64 `json:"projectId"`
}

//  获取待归库的尺寸列表请求参数
type GetToBeEnterRepositorySpecificationListInput struct {
	//  项目id
	ProjectID int64 `json:"projectId"`
}

type GraphDesc struct {
	//  接口错码说明
	Title string `json:"title"`
	//  详细说明
	Desc string `json:"desc"`
	//  错码列表
	ErrCodes []*ErrCodes `json:"errCodes"`
}

//  安装型钢需要的参数
type InstallLocationInput struct {
	//  安装编码
	LocationCode int64 `json:"locationCode"`
	//  识别码
	Identifier string `json:"identifier"`
}

//  安装码是否可用请求参数
type IsAccessLocationCodeInput struct {
	//  型钢识别码
	Identifier string `json:"identifier"`
	//  安装编码
	LocationCode int64 `json:"locationCode"`
}

type LoginRes struct {
	//  授权token
	AccessToken string `json:"accessToken"`
	//  过期时间戳(秒 7天)
	Expired int64 `json:"expired"`
	//  角色标识
	Role roles.GraphqlRole `json:"role"`
	//  角色名
	RoleName string `json:"roleName"`
}

//  分页参数
type PaginationInput struct {
	//  每页数量
	PageSize int64 `json:"pageSize"`
	//  指定哪个分页
	Page int64 `json:"page"`
	//  指定规格
	SpecificationID *int64 `json:"specificationId"`
	//  指定仓库
	RepositoryID *int64 `json:"repositoryId"`
}

type ProjectOrder2WorkshopDetailInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  规格id
	SpecificationID *int64 `json:"specificationId"`
	//  订单id
	OrderID int64 `json:"orderId"`
}

type ProjectOrder2WorkshopInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  订单id
	OrderID int64 `json:"orderId"`
	//  物流公司id
	ExpressCompanyID int64 `json:"expressCompanyId"`
	//  物流编号
	ExpressNo string `json:"expressNo"`
}

// " 待修改武钢详细信息请求参数
type ProjectSteel2BeChangeInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  规格尺寸id
	SpecificationID *int64 `json:"specificationId"`
	//  状态
	State *int64 `json:"state"`
}

type RepositoryLeaderItem struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Wechat string `json:"wechat"`
}

//  批量维修出库参数
type SetBatchOfMaintenanceSteelInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  维修厂 id
	MaintenanceID int64 `json:"maintenanceId"`
}

type SetBatchOfRepositorySteelInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  规格id
	SpecificationID int64 `json:"specificationId"`
	//  材料商id
	MaterialManufacturersID int64 `json:"materialManufacturersId"`
	//  生产商id
	ManufacturerID int64 `json:"manufacturerId"`
	//  生产日期
	ProducedAt time.Time `json:"producedAt"`
}

//  批量报废型钢参数
type SetBatchOfRepositorySteelScrapInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
}

//  设置密码的参数
type SetPasswordInput struct {
	Password string `json:"password"`
}

//  型钢归库参数
type SetProjectSteelEnterRepositoryInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  项目id
	ProjectID int64 `json:"projectId"`
}

//  修改项目型钢需要的参数
type SetProjectSteelInput struct {
	//  标识码列表
	IdentifierList []string `json:"identifierList"`
	//  型钢状态
	State int64 `json:"state"`
}

//  型钢出场请求参数
type SetProjectSteelOutOfWorkshopInput struct {
	//  识别码列表
	IdentifierList []string `json:"identifierList"`
	//  项目id
	ProjectID int64 `json:"projectId"`
}

type SetSteelIntoWorkshopInput struct {
	//  型钢标识码列表
	IdentifierList []string `json:"identifierList"`
	//  订单ID
	OrderID int64 `json:"orderId"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetOrderDetailInput struct {
	//  订单id
	ID int64 `json:"id"`
}

//  角色
type CreateInputUserRole string

const (
	//  仓库管理员
	CreateInputUserRoleRepositoryAdmin CreateInputUserRole = "repositoryAdmin"
	//  项目管理员
	CreateInputUserRoleProjectAdmin CreateInputUserRole = "projectAdmin"
	//  维修管理员
	CreateInputUserRoleMaintenanceAdmin CreateInputUserRole = "maintenanceAdmin"
)

var AllCreateInputUserRole = []CreateInputUserRole{
	CreateInputUserRoleRepositoryAdmin,
	CreateInputUserRoleProjectAdmin,
	CreateInputUserRoleMaintenanceAdmin,
}

func (e CreateInputUserRole) IsValid() bool {
	switch e {
	case CreateInputUserRoleRepositoryAdmin, CreateInputUserRoleProjectAdmin, CreateInputUserRoleMaintenanceAdmin:
		return true
	}
	return false
}

func (e CreateInputUserRole) String() string {
	return string(e)
}

func (e *CreateInputUserRole) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = CreateInputUserRole(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid CreateInputUserRole", str)
	}
	return nil
}

func (e CreateInputUserRole) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GetOrderListInputType string

const (
	//  确认订单
	GetOrderListInputTypeConfirmOrder GetOrderListInputType = "confirmOrder"
	//  待确认订单
	GetOrderListInputTypeToBeConfirm GetOrderListInputType = "toBeConfirm"
)

var AllGetOrderListInputType = []GetOrderListInputType{
	GetOrderListInputTypeConfirmOrder,
	GetOrderListInputTypeToBeConfirm,
}

func (e GetOrderListInputType) IsValid() bool {
	switch e {
	case GetOrderListInputTypeConfirmOrder, GetOrderListInputTypeToBeConfirm:
		return true
	}
	return false
}

func (e GetOrderListInputType) String() string {
	return string(e)
}

func (e *GetOrderListInputType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GetOrderListInputType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GetOrderListInputType", str)
	}
	return nil
}

func (e GetOrderListInputType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
