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
	seeds := []seed.Seed{}
	seeds = append( seeds, configsSeeders... )
	seeds = append( seeds, UsersSeeders... )
	seeds = append( seeds, rolesSeeders... )
	seeds = append(seeds, specificationinfoSeeds...)
	seeds = append(seeds, codeInfoSeeds...)
	seeds = append(seeds, fileSeeders...)
	seeds = append(seeds, repositorySeeder...)
	seeds = append(seeds, deviceSeeders...)

	return seeds
}


