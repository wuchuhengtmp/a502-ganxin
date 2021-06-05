/**
 * @Desc    The requests is part of http-api
 * @Author  wuchuheng<wuchuheng@163.com>
 * @Blog    https://wuchuheng.com
 * @wechat  wc20030318
 * @DATE    2021/6/4
 * @Listen  MIT
 */
package requests

import (
	"context"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"http-api/app/http/graph/model"
)

func ValidateCreateSpecificationRequest(ctx context.Context, input model.CreateSpecificationInput)  error {
	rules := govalidator.MapData{
		"length": []string{"isGreaterZero"},
		"weight": []string{"isGreaterZero"},
		"type":   []string{"min:6"},
	}
	message := govalidator.MapData{
		 "type": []string{
		 	"min:类型不能少于6个字",
		 },
	}
	opts := govalidator.Options{
		Data: &input,
		Rules: rules,
		TagIdentifier: "json",
		Messages: message,
	}
	res := govalidator.New(opts).ValidateStruct()
	if len(res) > 0 {
		for _, fieldErrors := range res {
			for _, err := range fieldErrors {
				return fmt.Errorf("%s", err)
			}
		}
	}

	return nil
}