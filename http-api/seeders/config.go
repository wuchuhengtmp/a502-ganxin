/**
 * @Desc    声明种子数据
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @DATE    2021/4/27
 * @Listen  MIT
 */
package seeders

import (
	"http-api/pkg/seed"
)

func All() []seed.Seed {
	return append(
		[]seed.Seed{},
		configsSeeders...
	)
}

