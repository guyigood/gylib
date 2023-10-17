package token

import (
	"gylib/common"
	"gylib/common/rediscomm"
	"sync"
)

type AccessToken struct {
	Slock    sync.Mutex
	Token    string
	WeixinID string
	Data     map[string]interface{}
	TimeOut  int
	Redis    *rediscomm.RedisComm
}

func NewAccessToken() *AccessToken {
	this := new(AccessToken)
	this.Slock.Lock()
	this.Data = make(map[string]interface{})
	this.Slock.Unlock()
	this.Redis = rediscomm.NewRedisComm()
	this.TimeOut = 7200
	return this
}

func (this *AccessToken) HasToken() bool {
	return this.Redis.SetKey(this.Token).HasKey()
}

func (this *AccessToken) SetToken(token string) *AccessToken {
	this.Token = token
	return this
}

func (this *AccessToken) SetTimeOut(timeout int) *AccessToken {
	this.TimeOut = timeout
	return this
}

func (this *AccessToken) GetTokenData() *AccessToken {
	data := this.Redis.SetKey(this.Token).Get_value()
	this.Slock.Lock()
	if data != nil {
		list, ok := data.(map[string]interface{})
		if !ok {
			this.Data = nil
		} else {
			this.Data = list
		}
	} else {
		this.Data = nil
	}
	this.Slock.Unlock()
	return this
}

func (this *AccessToken) DelToken() {
	this.Redis.SetKey(this.Token).DelKey()
}

func (that *AccessToken) GetTokenRawStr() interface{} {
	data := that.Redis.SetKey(that.Token).GetRawValue()
	if data == nil {
		return nil
	}
	return data
}

func (this *AccessToken) GetTokenValue(key string) interface{} {
	data := this.Redis.SetKey(this.Token).Get_value()
	if data == nil {
		return nil
	}
	this.Slock.Lock()
	defer this.Slock.Unlock()
	list, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	this.Data = list
	if common.Has_map_index(key, this.Data) {
		return this.Data[key]
	} else {
		return nil
	}
}

func (this *AccessToken) UpdateToken() {
	if this.Redis.SetKey(this.Token).HasKey() {
		raw := this.Redis.SetKey(this.Token).GetRawValue()
		this.Redis.SetKey(this.Token).SetExec("SETEX").SetTime(this.TimeOut).SetData(raw).SetRawValue()
	}
}

func (this *AccessToken) SetData(data map[string]interface{}) {
	this.Slock.Lock()
	defer this.Slock.Unlock()
	this.Data = data
	this.Redis.SetKey(this.Token).SetExec("SETEX").SetTime(this.TimeOut).SetData(this.Data).Set_value()
}

func (this *AccessToken) SetList(data map[string]interface{}) {
	this.Slock.Lock()
	defer this.Slock.Unlock()
	this.Data = data
	this.Redis.SetKey(this.Token).SetData(this.Data).SetList()
}

func (this *AccessToken) Save_Redis_Token() {
	if this.Token == "" {
		return
	}
	this.Slock.Lock()
	defer this.Slock.Unlock()
	this.Redis.SetKey(this.Token)
	this.Redis.Data = map[string]interface{}{"uuid": this.Token}
	this.Redis.Common_exec = "SETEX"
	this.Redis.Timeout = this.TimeOut
	this.Redis.Set_value()
}

func (this *AccessToken) GetNewToken() string {
	uu_id := common.Get_UUID()
	this.Token = uu_id
	return uu_id
}
