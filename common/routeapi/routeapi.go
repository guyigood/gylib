package routeapi

import (
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/rediscomm"
	"github.com/guyigood/gylib/common/webclient"
	"strings"
)

/*
	  调用示例
	    wb:=routeapi.NewRouteAPI()
		wb.Url=pubinit.WebData["api_url"]
		wb.Token_url=pubinit.WebData["token_url"]
		wb.Login_url=pubinit.WebData["login_url"]
		wb.User=pubinit.WebData["api_code"]
		wb.Pass=pubinit.WebData["api_pass"]
		wb.GetToken()
		wb.Get_Login_Token()
		data:=wb.SetData(this.Postdata).Post_route(pubinit.WebData["report_url"])
*/
type RouteApi struct {
	Wbc                       *webclient.Http_Client
	User                      string
	Pass                      string
	Url, Token_url, Login_url string
	token_key, Access_token   string
	Is_login                  bool
	Data                      map[string]interface{}
}

func NewRouteAPI() *RouteApi {
	this := new(RouteApi)
	this.Wbc = webclient.NewHttpClient()
	this.token_key = "webapi_access_token"
	this.Is_login = false
	this.S_init()
	return this
}

func (this *RouteApi) S_init() {
	this.Data = make(map[string]interface{})
}

func (this *RouteApi) GetToken() {
	if this.Access_token != "" {
		return
	}
	client := rediscomm.NewRedisComm()
	if client.SetKey(this.token_key).HasKey() {
		token_raw := client.SetKey(this.token_key).GetRawValue()
		if token_raw == nil {
			return
		} else {
			this.Access_token = datatype.Type2str(token_raw)
			return
		}
	} else {
		data := this.Get_client(this.Token_url)
		if data == nil {
			return
		} else {
			//fmt.Println(data)
			client.SetKey(this.token_key).SetData(data["data"]).SetExec("SETEX").SetTime(3600).SetRawValue()
			this.Access_token = datatype.Type2str(data["data"])
			return
		}

	}

}

func (this *RouteApi) Get_Login_Token() {
	this.GetToken()
	if this.Is_login {
		return
	}
	//redis := rediscomm.NewRedisComm()
	s_data := make(map[string]interface{})
	s_data["code"] = this.User
	s_data["pass"] = this.Pass
	s_data["access_token"] = this.Access_token
	result := this.Wbc.Https_post(this.Url+this.Login_url, s_data)
	list := datatype.String2Json(result)
	if list != nil {
		this.Is_login = true
		//this.Access_token = datatype.Type2str(list["data"])
		//redis.SetKey(this.token_key).SetData(this.Access_token).SetExec("SETEX").SetTime(3000).SetRawValue()
	}

}

func (this *RouteApi) SetData(data map[string]interface{}) *RouteApi {
	this.Data = data
	return this
}

func (this *RouteApi) Post_route(url string) map[string]interface{} {
	//this.Get_Login_Token()
	this.Data["access_token"] = this.Access_token
	data := this.Wbc.Https_post(this.Url+url, this.Data)
	//fmt.Println("post",data)
	if data == "" {
		return nil
	} else {
		return datatype.String2Json(data)
	}
}

func (this *RouteApi) Post_route_data(url string) []map[string]interface{} {
	//this.Get_Login_Token()
	this.Data["access_token"] = this.Access_token
	data := this.Wbc.Https_post(this.Url+url, this.Data)
	if data == "" {
		return nil
	} else {
		web_data := datatype.String2Json(data)
		result_data := make([]map[string]interface{}, 0)
		if datatype.Type2str(web_data["status"]) != "100" {
			return nil
		}
		list, ok := web_data["data"].([]interface{})
		if !ok {
			return nil
		}
		for _, v := range list {
			tmp, ok := v.(map[string]interface{})
			if ok {
				result_data = append(result_data, tmp)
			}
			//fmt.Println(tmp)
		}
		return result_data
	}
}

func (this *RouteApi) Get_route(params string) map[string]interface{} {
	//this.Get_Login_Token()

	if strings.Contains(params, "?") {
		params += "&access_token=" + this.Access_token
	} else {
		params = "?access_token=" + this.Access_token
	}
	//fmt.Println(this.Url+params)
	data := this.Wbc.HttpGet(this.Url + params)
	//fmt.Println(this.Url + params)
	if data == "" {
		return nil
	} else {
		return datatype.String2Json(data)
	}
}

func (this *RouteApi) Get_client(url string) map[string]interface{} {
	data := this.Wbc.HttpGet(this.Url + url)
	//fmt.Println(this.Url+url)
	if data == "" {
		return nil
	} else {
		return datatype.String2Json(data)
	}
}
