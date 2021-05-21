package requests

import (
	"github.com/thedevsaddam/govalidator"
	"http-api/app/models/article"
)

//
func ValidateArticleForm(data article.Article) map[string][]string {
	message := govalidator.MapData{
		"title": []string {
			"required:标题必填项",
			"min:标题长度需大于3",
			"max:标题长度度小于40",
		},
		"body": []string {
			"required:文章内容为必填项",
			"min:长度需大于 10",
		},
	}
	rules := govalidator.MapData{
		"title": []string{"required", "min:3", "max:40"},
		"body": []string{"required", "min:10"},
	}
	opts := govalidator.Options{
		Data: 			&data,
		Rules: 			rules,
		TagIdentifier:  "valid",
		Messages: 		message,
	}

	return govalidator.New(opts).ValidateStruct()
}
