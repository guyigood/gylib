package weixinmp

import (
	"encoding/json"
	"gylib/common/datatype"
	"gylib/common/redispack"
	"io/ioutil"
	"net/http"
	"time"
)

type AccessToken struct {
	AppId     string
	AppSecret string
	TmpName   string
	LckName   string
	Is_Redis  int
}

// get fresh access_token string
func (this *AccessToken) Fresh() (string, error) {
	redis_pool := redispack.Get_redis_pool()
	redis := redis_pool.Get()
	defer redis.Close()
	access_token, err := redis.Do("GET", redispack.Redis_data["redis_perfix"]+this.AppId+"_access_token")

	//fmt.Println("get_redis",datatype.Type2str(access_token),redispack.Redis_data["redis_perfix"]+this.AppId+"_access_token")
	if datatype.Type2str(access_token) == "" {
		if this.Is_Redis == 1 { //只允许从redis中读取
			error_ct := 0
			r_token := ""
			for {
				access_token1, _ := redis.Do("GET", redispack.Redis_data["redis_perfix"]+this.AppId+"_access_token")
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
			return r_token, nil
		}
		url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + this.AppId + "&secret=" + this.AppSecret
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// handle error
			return "", err
		}
		http_result := make(map[string]interface{})
		json.Unmarshal([]byte(body), &http_result)
		v, ok := http_result["access_token"]
		if ok {
			access_token = datatype.Type2str(v)
		} else {
			return "", err
		}
		redis.Do("SETEX", redispack.Redis_data["redis_perfix"]+this.AppId+"_access_token", 3600, access_token)
	}

	return datatype.Type2str(access_token), err
}

func (this *AccessToken) Get_Jsapi_ticket() (string, error) {
	redis_pool := redispack.Get_redis_pool()
	redis := redis_pool.Get()
	defer redis.Close()
	access_token, err := redis.Do("GET", redispack.Redis_data["redis_perfix"]+this.AppId+"_jsapi_ticket")
	if datatype.Type2str(access_token) == "" {
		accessToken, err := this.Fresh()
		if err != nil {
			return "", err
		}
		url := "https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=" + accessToken
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// handle error
			return "", err
		}
		http_result := make(map[string]interface{})
		json.Unmarshal([]byte(body), &http_result)
		v, ok := http_result["ticket"]
		if ok {
			access_token = datatype.Type2str(v)
		} else {
			return "", err
		}
		redis.Do("SETEX", redispack.Redis_data["redis_perfix"]+this.AppId+"_jsapi_ticket", 3600, access_token)
	}
	return datatype.Type2str(access_token), err
}
