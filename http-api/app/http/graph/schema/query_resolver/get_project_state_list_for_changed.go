/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/29
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/models/steels"
)

func (*QueryResolver) GetProjectStateListForChanged(ctx context.Context) (res []*steels.StateItem, err error) {

	for _, item := range []int64{
		steels.StateProjectWillBeUsed,    //【项目】-待使用
		steels.StateProjectInUse,         //【项目】-使用中
		steels.StateProjectException,     //【项目】-异常
		steels.StateProjectIdle,          //【项目】-闲置
		steels.StateProjectWillBeStore,   //【项目】-准备归库
		steels.StateProjectOnTheStoreWay, //【项目】-归库途中
	} {
		tmp := steels.StateItem{
			State: item,
			Desc:  steels.StateCodeMapDes[item],
		}
		res = append(res, &tmp)
	}

	return
}
