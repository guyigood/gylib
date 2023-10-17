package weixinsdk

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gylib/common/datatype"
	"gylib/common/redispack"
	"gylib/common/webclient"
	"io/ioutil"
	"strings"
	"time"
)

type Wxsdk struct {
	Access_token  string
	Appid, Appkey string
	RePool        *redis.Pool
	Is_Redis      int
	PostResult    string
}

func NewWxsdk() *Wxsdk {
	this := new(Wxsdk)
	this.Is_Redis = 0
	this.RePool = redispack.Get_redis_pool()
	return this
}

func (this *Wxsdk) DelAccessToken() {
	key := redispack.Redis_data["redis_perfix"] + this.Appid + "_access_token"
	redis := this.RePool.Get()
	defer redis.Close()
	redis.Do("DEL", key)

}

func (this *Wxsdk) Get_access_token() string {
	//redis_pool := redispack.Get_redis_pool()
	client := webclient.NewHttpClient()

	redis := this.RePool.Get()
	defer redis.Close()
	key := redispack.Redis_data["redis_perfix"] + this.Appid + "_access_token"
	access_token, err := redis.Do("GET", key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//fmt.Println("get_token",access_token,this)
	if datatype.Type2str(access_token) == "" || access_token == nil {

		if this.Is_Redis == 1 { //只允许从redis中读取
			error_ct := 0
			r_token := ""
			for {
				access_token1, err1 := redis.Do("GET", key)
				if err1 != nil {
					fmt.Println(err1)
					error_ct++
					if error_ct > 3 {
						break
					}
					continue
				}
				r_token = datatype.Type2str(access_token1)
				if r_token != "" {
					break
				}
				time.Sleep(time.Second * 1)
				error_ct++
				if error_ct > 3 {
					break
				}
			}
			return r_token
		}
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
		redis.Do("SETEX", key, 7000, access_token)
	}
	return datatype.Type2str(access_token)
}

func (this *Wxsdk) Send_wx_template(memo string) int {
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
		fmt.Println(string(body))
		fmt.Println(url)
	}

	return 2

}

func (this *Wxsdk) Send_New_wx_template(memo string) int {
	client := webclient.NewHttpClient()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + this.Access_token
	//data := fmt.Sprintf("{\"touser\":\"%s\",\"msgtype\":\"text\",\"text\":{\"content\":\"%s\"}}", userid, memo)
	resp, err := client.Client.Post(url, "application/json", strings.NewReader(memo))
	if err != nil {

		fmt.Println(resp)
		fmt.Println(err)
		return -1
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(memo)
		fmt.Println(string(body))
		fmt.Println(err)
		// handle error
		return -1
	}
	fmt.Println(string(body))
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

func (this *Wxsdk) Get_Jsapi_ticket() string {
	ticket := ""
	//redis_pool := redispack.Get_redis_pool()
	redis := this.RePool.Get()
	defer redis.Close()
	key := redispack.Redis_data["redis_perfix"] + this.Appid + "_jsapi_ticket"
	access_token, err := redis.Do("GET", key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//fmt.Println("redis",access_token,redispack.Redis_data["redis_perfix"]+this.Appid+"_jsapi_ticket")
	if access_token == nil {
		ticket = this.Get_Http_jsapi_ticket()
		if ticket != "" {
			redis.Do("SETEX", key, 7000, ticket)
		}
	} else {
		if datatype.Type2str(access_token) == "" {
			ticket = this.Get_Http_jsapi_ticket()
			if ticket != "" {
				redis.Do("SETEX", key, 7000, ticket)
			}
		} else {
			ticket = datatype.Type2str(access_token)
		}

	}
	return ticket
}

func (this *Wxsdk) Get_Http_jsapi_ticket() string {
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
