/**
 * @Desc    获取日志列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/8
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/logs"
	"http-api/app/models/users"
	"http-api/pkg/model"
)

func (*QueryResolver) GetLogList(ctx context.Context, input graphModel.GetLogListInput) (*logs.GetLogListRes, error) {
	if err := requests.ValidateGetLogListRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	logItem := logs.Logos{}
	userTable := users.Users{}.TableName()
	me := auth.GetUser(ctx)
	modelIns := model.DB.Debug().Model(&logItem).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.uid", userTable, userTable, logItem.TableName())).
		Where(fmt.Sprintf("%s.company_id = ?", userTable), me.CompanyId).
		Order(fmt.Sprintf("%s.id desc", logItem.TableName()))
	// 操作类型过滤
	if input.Type != nil {
		modelIns = modelIns.Where(fmt.Sprintf("%s.type = ?", logItem.TableName()), *input.Type)
	}

	res := logs.GetLogListRes{}
	err := modelIns.Count(&res.Total).Error
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	// 分页查询
	if !input.IsShowAll {
		modelIns = modelIns.Limit(int(*input.PageSize)).Offset(int((*input.Page - 1) * *input.PageSize))
	}
	if err := modelIns.Select(fmt.Sprintf("%s.*", logItem.TableName())).Find(&res.List).Error; err != nil {
		return nil, errors.ServerErr(ctx, err)
	}

	return &res, nil
}
