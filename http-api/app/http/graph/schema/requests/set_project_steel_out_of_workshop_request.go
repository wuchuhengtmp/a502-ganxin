/**
 * @Desc    型钢出场请求验证器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"sync"
)

func ValidateSetProjectSteelOutOfWorkshopRequest(ctx context.Context, input graphModel.SetProjectSteelOutOfWorkshopInput) error {
	// 检验有没有这个项目
	ps := StepsForProject{}
	if err := ps.CheckHasProject(ctx, input.ProjectID); err != nil {
		return err
	}
	steps := validateSetProjectSteelOutOfWorkshopRequestSteps{}
	var errors []error
	errChan := make(chan error, 10)
	errWg := &sync.WaitGroup{}
	limiter := make(chan bool, 10)
	go func() {
		errWg.Add(1)
		defer errWg.Done()
		for err := range errChan {
			errors = append(errors, err)
		}
	}()
	wg := &sync.WaitGroup{}
	for _, identifier := range input.IdentifierList {
		wg.Add(1)
		limiter <- true
		go steps.CheckIdentifier(ctx, identifier, wg, &errChan, limiter)
	}
	wg.Wait()
	// 结束错误收集
	close(errChan)
	errWg.Wait()
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

type validateSetProjectSteelOutOfWorkshopRequestSteps struct{}

func (*validateSetProjectSteelOutOfWorkshopRequestSteps) CheckIdentifier(
	ctx context.Context,
	identifier string,
	wg *sync.WaitGroup,
	errChan *chan error,
	limiter chan bool,
) {
	defer func() {
		wg.Done()
		<-limiter // 释放一个并发限制
	}()
	steps := StepsForProject{}
	// 检验有没有这根型钢
	if err := steps.CheckHasSteel(ctx, identifier); err != nil {
		*errChan <- err
	}
	// 检验型钢是否属于我
	if err := steps.CheckIsBelongMeByIdentifier(ctx, identifier); err != nil {
		*errChan <- err
	}
	// 检验型钢状态
	orderSpecificationSteelItem := order_specification_steel.OrderSpecificationSteel{}
	steelTable := steels.Steels{}.TableName()
	model.DB.Model(&orderSpecificationSteelItem).
		Joins(fmt.Sprintf("join %s ON %s.order_specification_steel_id = %s.id", steelTable, steelTable, orderSpecificationSteelItem.TableName())).
		Where(fmt.Sprintf("%s.identifier = ?", steelTable), identifier).
		First(&orderSpecificationSteelItem)
	isAllow := false
	for _, state := range []int64 {
		steels.StateProjectWillBeUsed,    //【项目】-待使用
		steels.StateProjectInUse,         //【项目】-使用中
		steels.StateProjectException,     //【项目】-异常
		steels.StateProjectIdle,          //【项目】-闲置
		steels.StateProjectWillBeStore,   //【项目】-准备归库
	} {
		if orderSpecificationSteelItem.State == state {
			isAllow = true
			break
		}
	}
	if !isAllow {
		*errChan <- fmt.Errorf( "型钢标识码为:%s的状态为%s,不能出场",  identifier, steels.StateCodeMapDes[orderSpecificationSteelItem.State])
	}
}
