/**
 * @Desc    The query_resolver is part of http-api
 * @Author  wuchuheng<root@wuchuheng.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/7/14
 * @Listen  MIT
 */
package requests

import (
	"fmt"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/files"
	"http-api/pkg/model"
)

func ValidateSetCompanyInfoRequest(input graphModel.SetCompanyInfoInput) error  {
	fileItem := files.File{}
	err := model.DB.Model(&files.File{}).Where("id = ?", input.TutorFileID).First(&fileItem).Error
	if err != nil && err.Error() == "record not found" {
		return fmt.Errorf("文件id为:%d 不存在", input.TutorFileID)
	}

	return err
}
