/**
 * @Desc    型钢入库解析器
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/11
 * @Listen  MIT
 */
package mutation_resolver

import (
	"context"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/steels"
)

func (*MutationResolver) CreateSteel(ctx context.Context, input model.CreateSteelInput) ([]*steels.Steels, error) {
	if err := requests.ValidateCreateSteelRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	me := auth.GetUser(ctx)
	var ss []*steels.Steels
	for _, identifier := range input.IdentifierList {
		s := steels.Steels{
			Identifier:             identifier,
			CreatedUid:             me.ID,
			State:                  steels.StateInStore,
			SpecificationId:        input.SpecificationID,
			CompanyId:              me.CompanyId,
			RepositoryId:           input.RepositoryID,
			MaterialManufacturerId: input.MaterialManufacturerID,
			ManufacturerId:         input.ManufacturerID,
			ProducedDate:           input.ProducedDate,
		}
		ss = append(ss, &s)
	}
	steelModel := steels.Steels{}
	if err := steelModel.CreateMultipleSteel(ctx, ss); err != nil {
		return ss, errors.ServerErr(ctx, err)
	}

	return ss, nil
}
