package controller

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"

	"mask_api_gin/src/internal/config"
	"mask_api_gin/src/internal/infra/cache"
	"mask_api_gin/src/internal/shared/constants"
	"mask_api_gin/src/internal/shared/generate"
	"mask_api_gin/src/internal/shared/logger"
	"mask_api_gin/src/internal/shared/parse"
	"mask_api_gin/src/internal/shared/regular"
	"mask_api_gin/src/internal/web/resp"
	systemService "mask_api_gin/src/modules/system/service"
)

// NewCaptcha 实例化控制层
var NewCaptcha = &CaptchaController{
	sysConfigService: systemService.NewSysConfig,
}

// CaptchaController 验证码操作 控制层处理
//
// PATH /api/v1/auth
type CaptchaController struct {
	sysConfigService *systemService.SysConfig // 参数配置服务
}

// Image 获取验证码图片
//
// GET /captcha-image
func (s CaptchaController) Image(c *gin.Context) {
	// 从 sys_config 获取验证码开关：true开启，false关闭
	captchaEnabledStr := s.sysConfigService.FindValueByKey("sys.account.captchaEnabled")
	captchaEnabled := parse.Boolean(captchaEnabledStr)
	if !captchaEnabled {
		c.JSON(200, resp.OkData(map[string]any{
			"enabled": captchaEnabled,
		}))
		return
	}

	// 生成唯一标识
	verifyKey := ""
	data := map[string]any{
		"enabled": captchaEnabled,
		"uuid":    "",
		"img":     "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7",
	}

	// 验证码有效期，单位秒
	captchaExpiration := 2 * 60 * time.Second
	// 从 sys_config 获取验证码类型：math 数值计算，char 字符验证
	captchaTypeStr := s.sysConfigService.FindValueByKey("sys.account.captchaType")
	if captchaTypeStr == "math" {
		// captchaMathConfig 数学验证码视觉参数（sys.account.captchaConfigMath，JSON 类型）。
		// 缺失或解析失败时各字段为零值，由调用方兜底。
		var mc struct {
			Height     int    `json:"height"`
			Width      int    `json:"width"`
			Noise      int    `json:"noise"`
			Color      bool   `json:"color"`
			Background string `json:"background"`
		}
		rawStr := s.sysConfigService.FindValueByKey("sys.account.captchaConfigMath")
		if rawStr != "" {
			if err := json.Unmarshal([]byte(rawStr), &mc); err != nil {
				logger.Errorf("captcha math config unmarshal err: %v", err)
			}
		}
		driverCaptcha := &base64Captcha.DriverMath{
			//Height png height in pixel.
			Height: mc.Height,
			// Width Captcha png width in pixel.
			Width: mc.Width,
			//NoiseCount text noise count.
			NoiseCount: mc.Noise,
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if mc.Background != "" {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(mc.Background)
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("generate id question answer %s %s : %v", captchaTypeStr, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			verifyKey = constants.CACHE_CAPTCHA_CODE + ":" + id
			_ = cache.Set("", verifyKey, answer, captchaExpiration)
		}
	}
	if captchaTypeStr == "char" {
		// captchaCharConfig 字符验证码视觉参数（sys.account.captchaConfigChar，JSON 类型）。
		var cc struct {
			Height     int    `json:"height"`
			Width      int    `json:"width"`
			Noise      int    `json:"noise"`
			Size       int    `json:"size"`
			Color      bool   `json:"color"`
			Background string `json:"background"`
			Chars      string `json:"chars"`
		}
		rawStr := s.sysConfigService.FindValueByKey("sys.account.captchaConfigChar")
		if rawStr != "" {
			if err := json.Unmarshal([]byte(rawStr), &cc); err != nil {
				logger.Errorf("captcha char config unmarshal err: %v", err)
			}
		}
		driverCaptcha := &base64Captcha.DriverString{
			//Height png height in pixel.
			Height: cc.Height,
			// Width Captcha png width in pixel.
			Width: cc.Width,
			//NoiseCount text noise count.
			NoiseCount: cc.Noise,
			//Length random string length.
			Length: cc.Size,
			//Source is a unicode which is the rand string from.
			Source: cc.Chars,
			//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
			ShowLineOptions: base64Captcha.OptionShowHollowLine,
		}
		if cc.Background != "" {
			//BgColor captcha image background color (optional)
			driverCaptcha.BgColor = parse.Color(cc.Background)
		}
		// 验证码生成
		id, question, answer := driverCaptcha.GenerateIdQuestionAnswer()
		// 验证码表达式解析输出
		item, err := driverCaptcha.DrawCaptcha(question)
		if err != nil {
			logger.Infof("generate id question answer %s %s : %v", captchaTypeStr, question, err)
		} else {
			data["uuid"] = id
			data["img"] = item.EncodeB64string()
			verifyKey = constants.CACHE_CAPTCHA_CODE + ":" + id
			_ = cache.Set("", verifyKey, strings.ToLower(answer), captchaExpiration)
		}
	}

	// 本地开发下返回验证码结果，方便接口调试
	if config.Env() == "local" {
		text, _ := cache.Get("", verifyKey)
		data["text"] = text
		c.JSON(200, resp.OkData(data))
		return
	}
	c.JSON(200, resp.OkData(data))
}

// SMS 发送手机号登录验证码
//
// POST /captcha-sms
func (s CaptchaController) SMS(c *gin.Context) {
	var body struct {
		Phone string `json:"phone" binding:"required"` // 手机号码
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		errMsgs := fmt.Sprintf("bind err: %s", resp.FormatBindError(err))
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_PARSER, errMsgs))
		return
	}

	if !regular.ValidMobile(body.Phone) {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "请输入正确的手机号码"))
		return
	}
	if config.Env() != "local" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "短信服务未配置"))
		return
	}

	// 验证码过期时间，单位秒
	expiresInSeconds := parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaSMSExpiresIn"))
	// 验证码发送间隔，单位秒
	resendInterval := parse.Number(s.sysConfigService.FindValueByKey("sys.account.captchaSMSResendInterval"))
	if expiresInSeconds <= 0 {
		expiresInSeconds = 300
	}
	if resendInterval <= 0 {
		resendInterval = 60
	}

	cooldownKey := fmt.Sprintf("%s:%s", constants.CACHE_SMS_CODE_COOLDOWN, body.Phone)
	if cooldown, _ := cache.Get("", cooldownKey); cooldown != "" {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "验证码发送过于频繁，请稍后再试"))
		return
	}

	code := fmt.Sprintf("%06d", generate.Number(6))
	uuid := generate.String(32)
	codeKey := fmt.Sprintf("%s:%s", constants.CACHE_SMS_CODE, uuid)
	codeValue := fmt.Sprintf("%s:%s", body.Phone, code)
	if err := cache.Set("", codeKey, codeValue, time.Duration(expiresInSeconds)*time.Second); err != nil {
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "验证码生成失败"))
		return
	}
	if err := cache.Set("", cooldownKey, body.Phone, time.Duration(resendInterval)*time.Second); err != nil {
		_ = cache.Del("", codeKey)
		c.JSON(422, resp.CodeMsg(resp.CODE_PARAM_CHEACK, "验证码生成失败"))
		return
	}

	data := map[string]any{
		"uuid":            uuid,
		"expires_in":      expiresInSeconds,
		"resend_interval": resendInterval,
	}
	if config.Env() == "local" {
		data["text"] = code
	}
	c.JSON(200, resp.OkData(data))
}
