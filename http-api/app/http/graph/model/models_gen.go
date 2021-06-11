// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"http-api/app/models/roles"
	"io"
	"strconv"
	"time"
)

type CompanyItemRes struct {
	ID int64 `json:"id"`
	//  公司名
	Name string `json:"name"`
	//  用于型钢编码生成
	PinYin string `json:"pinYin"`
	//  APP 企业宗旨
	Symbol string `json:"symbol"`
	//  logo文件
	LogoFile *FileItem `json:"logoFile"`
	//  app背景文件
	BackgroundFile *FileItem `json:"backgroundFile"`
	//  账号状态
	IsAble bool `json:"isAble"`
	//  公司的电话
	Phone string `json:"phone"`
	//  公司的微信
	Wechat string `json:"wechat"`
	//  开始时间
	StartedAt time.Time `json:"startedAt"`
	//  结束时间
	EndedAt time.Time `json:"endedAt"`
	//  创建时间
	CreatedAt time.Time `json:"createdAt"`
	//  管理员名称
	AdminName string `json:"adminName"`
	//  管理员手机
	AdminPhone string `json:"adminPhone"`
	//  管理员微信
	AdminWechat string `json:"adminWechat"`
	//  管理员头像
	AdminAvatar *FileItem `json:"adminAvatar"`
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

type GraphDesc struct {
	//  接口错码说明
	Title string `json:"title"`
	//  详细说明
	Desc string `json:"desc"`
	//  错码列表
	ErrCodes []*ErrCodes `json:"errCodes"`
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

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
