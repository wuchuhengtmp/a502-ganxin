package mutation_resolver

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/http/graph/schema/requests"
	"http-api/app/models/configs"
	"time"
)

var msgKeyMapCode = make(map[string]struct {
	Code      string
	CreatedAt time.Time
})

func (*MutationResolver) CreateCode(ctx context.Context, input graphModel.GetCodeForForgetPasswordInput) (*graphModel.GetCodeForForgetPasswordRes, error) {
	if err := requests.ValidateCreateCodeRequest(ctx, input); err != nil {
		return nil, errors.ValidateErr(ctx, err)
	}
	accessKey := configs.GetGlobalVal(configs.SMS_ACCESS_KEY)
	accessSecretKey := configs.GetGlobalVal(configs.SMS_ACCESS_SECRET_KEY)
	client, err := dysmsapi.NewClientWithAccessKey("cn-shenzhen", accessKey, accessSecretKey)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = input.Phone
	sign := configs.GetGlobalVal(configs.SMS_SIGN)
	request.SignName = sign
	template := configs.GetGlobalVal(configs.SMS_TEMPLATECODE)
	request.TemplateCode = template
	code := fmt.Sprintf("%d", time.Now().Unix())[4:]
	request.TemplateParam = "{\"code\":" + code + "}"
	key := fmt.Sprintf("%d", time.Now().UnixNano())
	msgKeyMapCode[key] = struct {
		Code      string
		CreatedAt time.Time
	}{Code: code, CreatedAt: time.Now()}

	_, err = client.SendSms(request)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res := graphModel.GetCodeForForgetPasswordRes{
		Key: key,
	}

	return &res, nil
}
