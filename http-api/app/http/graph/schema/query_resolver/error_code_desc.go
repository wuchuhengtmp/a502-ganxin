/**
 * @Desc    错码码说明
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/5/31
 * @Listen  MIT
 */
package query_resolver

import (
	"context"
	"http-api/app/http/graph/errors"
	"http-api/app/http/graph/model"
)

func (*QueryResolver)ErrorCodeDesc(ctx context.Context) (*model.GraphDesc, error) {
	res := model.GraphDesc {
		Title: "这是接口错码码说明文档",
		Desc: `
		错码在错码列表中的extensions字段中如:
		{
		  "errors": [
			{
			  "message": "hello",
			  "path": [
				"errorCodeDesc"
			  ],
			  "extensions": {
				"code": 4100
			  }
			}
		  ],
		  "data": null
		}
		字段:code： 4100就是了
`,

		ErrCodes: []*model.ErrCodes{
			{
				Code: errors.ServerErrCode,
				Desc: "服务器端出了问，请联系管理员解决",
			},
			{
				Code: errors.InvalidErrCode,
				Desc: "验证错误码",
			},
			{
				Code: errors.AccessDenyErrCode,
				Desc: "权限限制错误码",
			},
			{
				Code: errors.DeviceDeniedCode,
				Desc: "设备禁用错误码",
			},
		},
	}

	return &res, nil
}
