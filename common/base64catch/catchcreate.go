package base64catch

import (
	"github.com/mojocn/base64Captcha"
)

type CodeCaptcha struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
	store         base64Captcha.Store
}

func NewCodecaptcha(param string) *CodeCaptcha {
	this := new(CodeCaptcha)
	this.CaptchaType = param
	this.store = base64Captcha.DefaultMemStore
	return this
}

func (this *CodeCaptcha) CodeCaptchaCreate() (string, string) {
	var driver base64Captcha.Driver
	var param CodeCaptcha
	//choose driver
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	c := base64Captcha.NewCaptcha(driver, this.store)
	id, b64s, _ := c.Generate()
	return id, b64s
}

func (this *CodeCaptcha) CaptchaVerify(id, verify string) bool {
	if this.store.Verify(id, verify, true) {
		return true
	}
	return false
}
