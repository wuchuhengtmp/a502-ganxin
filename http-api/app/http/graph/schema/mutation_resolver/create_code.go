package mutation_resolver

import (
	"context"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"http-api/app/http/graph/errors"
	graphModel "http-api/app/http/graph/model"
	"http-api/app/models/configs"
	"http-api/app/models/users"
	"http-api/pkg/model"
	"regexp"
	"time"
)

type SmsKeyMapCodeType map[string]struct {
	Code      string
	CreatedAt time.Time
	Phone     string
}
/**
 * 短信验证码缓存位置
 */
var SmsKeyMapCode = make(SmsKeyMapCodeType)

/**
 * 手机号最近获取到验证码的时间
 */
var PhoneMapLastTimeForSMS = make(map[string]time.Time)

func init()  {
	// 清理过期的验证码缓存数据
	go func() {
		ticker := time.NewTicker(time.Minute * 30)
		for {
			<- ticker.C
			for key, item := range SmsKeyMapCode {
				timeLen := time.Now().Unix() - item.CreatedAt.Unix()
				if timeLen > 60 * 30 {
					delete(SmsKeyMapCode, key)
				}
			}
		}
	}()
	// 清理没用的上次获取手机号的时间缓存数据
	go func() {
		secondLen := time.Duration( 60)
		ticker := time.NewTicker(time.Second * secondLen)
		for {
			<- ticker.C
			for key, item := range PhoneMapLastTimeForSMS {
				if item.Unix() + int64(secondLen) < time.Now().Unix() {
					delete(PhoneMapLastTimeForSMS, key)
				}
			}
		}
	}()

}

func (*MutationResolver) CreateCode(ctx context.Context, input graphModel.GetCodeForForgetPasswordInput) (*graphModel.GetCodeForForgetPasswordRes, error) {
	if err := ValidateCreateCodeRequest(input); err != nil {
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
	SmsKeyMapCode[key] = struct {
		Code      string
		CreatedAt time.Time
		Phone     string
	}{Code: code, CreatedAt: time.Now(), Phone: input.Phone}
	PhoneMapLastTimeForSMS[input.Phone] = time.Now()

	_, err = client.SendSms(request)
	if err != nil {
		return nil, errors.ServerErr(ctx, err)
	}
	res := graphModel.GetCodeForForgetPasswordRes{
		Key: key,
	}

	return &res, nil
}


func ValidateCreateCodeRequest(input graphModel.GetCodeForForgetPasswordInput) error {
	_, err := regexp.Match("^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199)\\d{8}$", []byte(input.Phone) )
	if err != nil {
		return fmt.Errorf("手机号为: %s 不是正确的手机号", input.Phone)
	}
	userItem := users.Users{}
	err = model.DB.Model(&userItem).Where("phone = ?", &input.Phone).Find(&userItem).Error
	if err != nil {
		return fmt.Errorf("没有手机号为: %s 的用户", input.Phone)
	}
	lastTime, ok := PhoneMapLastTimeForSMS[input.Phone]
	if ok {
		if time.Now().Unix() - lastTime.Unix() < 60{
			return fmt.Errorf("请在1分钟后，重新获取")
		}
	}

	return nil
}
