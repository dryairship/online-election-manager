package utils

import "github.com/mojocn/base64Captcha"

var captchaOptions = base64Captcha.ConfigDigit{
	CaptchaLen: 5,
	DotCount:   80,
	Height:     80,
	MaxSkew:    0.4,
	Width:      200,
}

func CreateCaptcha() (string, string) {
	captchaID, captcha := base64Captcha.GenerateCaptcha("", captchaOptions)
	base64 := base64Captcha.CaptchaWriteToBase64Encoding(captcha)
	return captchaID, base64
}

func VerifyCaptcha(id, value string) bool {
	return base64Captcha.VerifyCaptcha(id, value)
}
