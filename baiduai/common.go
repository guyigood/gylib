package baiduai

import (
	"fmt"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/rediscomm"
	"github.com/guyigood/gylib/common/webclient"
)

type BaiDuAI struct {
	AppId       string
	AppKey      string
	SecretKey   string
	AccessToken string
	imgBase64   string
	OcrResult   *OcrReturn
	IorcIn      *IocrInput
	curlweb     *webclient.Http_Client
}

func NewBaiDuAI(appid, appkey, seckey string) *BaiDuAI {
	this := new(BaiDuAI)
	this.AppId = appid
	this.AppKey = appkey
	this.SecretKey = seckey
	this.StructInit()
	return this
}

func (this *BaiDuAI) StructInit() {
	this.curlweb = webclient.NewHttpClient()
	this.curlweb.Init_HTTPClient()
}

func (this *BaiDuAI) GetAccessToken() {
	client := rediscomm.NewRedisComm()
	key := "bdai_" + this.AppId + "_token"
	if client.SetKey(key).HasKey() {
		this.AccessToken = datatype.Type2str(client.SetKey(key).GetRawValue())
		return
	}
	url := "https://aip.baidubce.com/oauth/2.0/token"
	data := make(map[string]interface{})
	data["grant_type"] = "client_credentials"
	data["client_id"] = this.AppKey
	data["client_secret"] = this.SecretKey
	result := this.curlweb.Https_post(url, data)
	r_data := datatype.String2Json(result)
	fmt.Println("token", result, r_data)
	if common.Has_map_index("access_token", r_data) {
		this.AccessToken = datatype.Type2str(r_data["access_token"])

		client.SetKey(key).SetData(this.AccessToken).SetExec("SETEX").SetTime(86400).SetRawValue()
	}
}
