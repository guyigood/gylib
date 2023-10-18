package redpack

import (
	"fmt"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/weixinsdk/wxpay"
)

type WxHongBao struct {
	AppId  string // 微信公众平台应用ID
	MchId  string // 微信支付商户平台商户号
	ApiKey string // 微信支付商户平台API密钥
	// 微信支付商户平台证书路径
	CertFile   string
	KeyFile    string
	RootcaFile string
	Wxuser     map[string]string
}

func NewWxHongBao(wxuserdata map[string]string) *WxHongBao {
	this := new(WxHongBao)
	this.Wxuser = make(map[string]string)
	this.Wxuser = wxuserdata
	this.AppId = this.Wxuser["appid"]
	this.ApiKey = this.Wxuser["apipass"]
	this.MchId = this.Wxuser["mchid"]
	this.CertFile = this.Wxuser["certfile"]
	this.KeyFile = this.Wxuser["keyfile"]
	this.RootcaFile = this.Wxuser["rootca"]
	return this
}

func (this *WxHongBao) Send_Pay(data map[string]interface{}) bool {
	c := wxpay.NewClient(this.AppId, this.MchId, this.ApiKey)
	// 附着商户证书
	err := c.WithCert(this.CertFile, this.KeyFile, this.RootcaFile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	params := make(wxpay.Params)
	for key, v := range data {
		switch v.(type) {
		case string:
			params.SetString(key, v.(string))
		case int64:
			params.SetInt64(key, v.(int64))
		}
	}
	// 查询企业付款接口请求参数
	params.SetString("mch_appid", c.AppId)
	params.SetString("mchid", c.MchId)
	params.SetString("nonce_str", common.RandomStr(32)) // 随机字符串
	params.SetString("sign", c.Sign(params))            // 签名

	// 查询企业付款接口请求URL
	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	ret, err := c.Post(url, params, true)

	if err != nil {
		return false
	} else {
		if ret["return_code"] == "SUCCESS" {
			return true
		} else {
			return false
		}
	}
}

func (this *WxHongBao) Send_Redpack(data map[string]interface{}) bool {
	c := wxpay.NewClient(this.AppId, this.MchId, this.ApiKey)
	// 附着商户证书
	err := c.WithCert(this.CertFile, this.KeyFile, this.RootcaFile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	params := make(wxpay.Params)
	for key, v := range data {
		switch v.(type) {
		case string:
			params.SetString(key, v.(string))
		case int64:
			params.SetInt64(key, v.(int64))
		}
	}
	// 查询企业付款接口请求参数
	params.SetString("mch_appid", c.AppId)
	params.SetString("mchid", c.MchId)
	params.SetString("nonce_str", common.RandomStr(32)) // 随机字符串
	params.SetString("sign", c.Sign(params))            // 签名

	// 查询企业付款接口请求URL
	url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
	ret, err := c.Post(url, params, true)

	if err != nil {
		return false
	} else {
		if ret["return_code"] == "SUCCESS" {
			return true
		} else {
			return false
		}
	}
}
