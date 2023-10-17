package weixinsend

import (
	"encoding/json"
	"fmt"
	"gylib/common/datatype"
	"gylib/common/rediscomm"
	"gylib/common/webclient"
	"io/ioutil"
	"strings"
)

type WxSendSdk struct {
	Access_token                          string
	Appid, Appkey, AToken_key, JToken_key string
	RePool                                *rediscomm.RedisComm
	PostResult                            string
}

func NewWxSendSdk(appid, appkey string) *WxSendSdk {
	this := new(WxSendSdk)
	this.RePool = rediscomm.NewRedisComm()
	this.Appid = appid
	this.Appkey = appkey
	this.AToken_key = this.Appid + "_access_token"
	this.JToken_key = this.Appid + "_jsapi_ticket"
	return this
}

func (this *WxSendSdk) DelAccessToken() {
	this.RePool.SetKey(this.AToken_key).DelKey()
}

func (this *WxSendSdk) Get_access_token() string {
	//redis_pool := redispack.Get_redis_pool()
	access_token := ""
	client := webclient.NewHttpClient()
	if this.RePool.SetKey(this.AToken_key).HasKey() {
		access_token = datatype.Type2str(this.RePool.SetKey(this.AToken_key).GetRawValue())
		this.Access_token = access_token
		return access_token
	}
	if access_token == "" {
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + this.Appid + "&secret=" + this.Appkey
		res, err := client.Client.Get(url)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			// handle error
			return ""
		}
		http_result := make(map[string]interface{})
		json.Unmarshal([]byte(body), &http_result)

		v, ok := http_result["access_token"].(string)
		if ok {
			access_token = v
		} else {
			return ""
		}
		this.Access_token = access_token
		this.RePool.SetKey(this.AToken_key).SetExec("SETEX").SetTime(7000).SetData(access_token).SetRawValue()
	}
	return access_token
}

func (this *WxSendSdk) Send_wx_template(memo string) int {
	client := webclient.NewHttpClient()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + this.Access_token
	//data := fmt.Sprintf("{\"touser\":\"%s\",\"msgtype\":\"text\",\"text\":{\"content\":\"%s\"}}", userid, memo)
	resp, err := client.Client.Post(url, "application/json", strings.NewReader(memo))
	if err != nil {
		fmt.Println(resp)
		fmt.Println(err, this.Access_token)
		return 0
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(string(body))
		fmt.Println(err)
		// handle error
		return 0
	}
	//fmt.Println(string(body))
	this.PostResult = string(body)
	aStr := make(map[string]interface{})
	json.Unmarshal([]byte(body), &aStr)
	v := datatype.Type2str(aStr["errcode"])
	if v == "0" {
		return 1
	} else {
		fmt.Println(url)
	}

	return 2

}

func (this *WxSendSdk) Get_Jsapi_ticket() string {
	ticket := ""
	//redis_pool := redispack.Get_redis_pool()
	if this.RePool.SetKey(this.JToken_key).HasKey() {
		ticket = datatype.Type2str(this.RePool.SetKey(this.JToken_key).GetRawValue())
		return ticket
	}
	if ticket == "" {
		ticket = this.Get_Http_jsapi_ticket()
		if ticket != "" {
			this.RePool.SetKey(this.JToken_key).SetExec("SETEX").SetTime(7000).SetData(ticket).SetRawValue()
		}
	}
	return ticket
}

func (this *WxSendSdk) Get_Http_jsapi_ticket() string {
	client := webclient.NewHttpClient()
	accessToken := this.Get_access_token()
	ticket := ""
	url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken
	res, err := client.Client.Get(url)
	//fmt.Println(url, accessToken, res)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		// handle error
		return ""
	}
	http_result := make(map[string]interface{})
	json.Unmarshal([]byte(body), &http_result)

	v, ok := http_result["ticket"]
	if ok {
		ticket = datatype.Type2str(v)
	} else {
		return ""
	}
	return ticket
}
