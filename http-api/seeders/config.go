/**
 * @Desc    声明种子数据
 * @Author  wuchuheng<root@wuchuheng.com>
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
	seeds = append(seeds, companySeeders...)
	seeds = append(seeds, projectSeeder...)
	seeds = append(seeds, steelsSeeds...)
	seeds = append(seeds, MaintenanceSeeder...)

	return seeds
}


