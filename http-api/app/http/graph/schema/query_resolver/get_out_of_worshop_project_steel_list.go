/**
 * @Desc    获取出场的项目列表解析器
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/28
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"fmt"
	"http-api/app/http/graph/auth"
	"http-api/app/http/graph/errors"
	"http-api/app/models/order_specification"
	"http-api/app/models/order_specification_steel"
	"http-api/app/models/orders"
	"http-api/app/models/project_leader"
	"http-api/app/models/projects"
	"http-api/app/models/steels"
	"http-api/pkg/model"
	"sync"
)

func (*QueryResolver) GetOutOfWorkshopProjectList(ctx context.Context) (res []*projects.Projects, err error) {
	var myLeadedProjects []projects.Projects
	var runtimeErrors []error
	runtimeErrorsChange := make(chan error, 10)
	me := auth.GetUser(ctx)
	projectTable := projects.Projects{}.TableName()
	projectLeaderTable := project_leader.ProjectLeader{}.TableName()
	err = model.DB.Model(&projects.Projects{}).
		Select(fmt.Sprintf("%s.*", projectTable)).
		Joins(fmt.Sprintf("join %s ON %s.project_id = %s.id", projectLeaderTable, projectLeaderTable, projectTable)).
		Where(fmt.Sprintf("%s.uid = ?", projectLeaderTable), me.Id).
		First(&myLeadedProjects).
		Error
	if err != nil {
		return
	}
	// 检验项目是否在场地的检验的限制并发量
	limiter := make(chan bool, 10)
	defer close(limiter)
	// 收集检验通过的项目
	wg := &sync.WaitGroup{}
	collectWg := &sync.WaitGroup{}
	allProjectInWorkshopChan := make(chan projects.Projects, 10) // 收集检验出属于场地的项目
	go func() {
		collectWg.Add(1)
		defer collectWg.Done()
		for projectInWorkshop := range allProjectInWorkshopChan {
			res = append(res, &projectInWorkshop)
		}
	}()
	// 收集协程中可能发生的错误
	collectErrWg := &sync.WaitGroup{}
	go func() {
		collectErrWg.Add(1)
		defer collectErrWg.Done()
		for err := range runtimeErrorsChange {
			runtimeErrors = append(runtimeErrors, err)
		}
	}()
	// 过虑我的项目，如果项目有型钢在场地中，则触发收集
	for _, projectItem := range myLeadedProjects {
		wg.Add(1)
		limiter <- true
		projectItem := projectItem
		go func() {
			defer wg.Done()
			steps := GetOutOfWorkshopProjectListSteps{}
			ok, err := steps.CheckIsStateOfWorkshop(projectItem)
			// 有错误则收集错误
			if err != nil {
				runtimeErrorsChange <- err
			} else {
				// 收集满足条件的项目
				if ok {
					allProjectInWorkshopChan <- projectItem
				}
			}
			<-limiter
		}()
	}
	wg.Wait()                       // 检验完毕
	close(allProjectInWorkshopChan) // 关闭收集合格项目管道
	collectWg.Wait()                // 收集结果完毕
	close(runtimeErrorsChange)      // 关闭错误收集管道
	collectErrWg.Wait()             // 错误收集完毕
	if len(runtimeErrors) > 0 {
		return res, errors.ValidateErr(ctx, runtimeErrors[0])
	}

	return
}

//获取出场的项目列表解决步骤
type GetOutOfWorkshopProjectListSteps struct{}

// 只要有型钢还留在场地中， 则判定项目还场地阶段
func (*GetOutOfWorkshopProjectListSteps) CheckIsStateOfWorkshop(projectItem projects.Projects) (bool, error) {
	orderSpecificationSteelTable := order_specification_steel.OrderSpecificationSteel{}.TableName()
	orderTable := orders.Order{}.TableName()
	orderSpecificationTable := order_specification.OrderSpecification{}.TableName()
	err := model.DB.Model(&order_specification_steel.OrderSpecificationSteel{}).
		Select(fmt.Sprintf("%s.*", orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_specification_id", orderSpecificationTable, orderSpecificationTable,orderSpecificationSteelTable)).
		Joins(fmt.Sprintf("join %s ON %s.id = %s.order_id", orderTable, orderTable, orderSpecificationTable)).
		Where(fmt.Sprintf("%s.state in ?", orderSpecificationSteelTable), steels.GetStateForProject()).
		Where(fmt.Sprintf("%s.project_id = ?",orderTable), projectItem.ID).
		First(&order_specification_steel.OrderSpecificationSteel{}).
		Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
