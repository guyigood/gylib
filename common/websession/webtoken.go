package websession

import (
	"encoding/json"
	"gylib/common"
	"gylib/common/datatype"
	"gylib/common/redisclient"
	"sync"
)

type WebAccessToken struct {
	Token   string
	Slock   sync.RWMutex
	Data    map[string]interface{}
	TimeOut int64
	Redis   *redisclient.RedisClient
}

func NewWebAccessToken() *WebAccessToken {
	that := new(WebAccessToken)
	that.Slock.Lock()
	that.Data = make(map[string]interface{})
	that.Slock.Unlock()
	that.Redis = redisclient.NewRedisCient()
	that.TimeOut = 7200
	return that
}

func (that *WebAccessToken) HasToken() bool {
	return that.Redis.SetKey(that.Token).HasKey()
}

func (that *WebAccessToken) SetToken(token string) *WebAccessToken {
	that.Token = token
	return that
}

func (that *WebAccessToken) SetTimeOut(timeout int64) *WebAccessToken {
	that.TimeOut = timeout
	return that
}

func (that *WebAccessToken) GetTokenData() *WebAccessToken {
	data := that.Redis.SetKey(that.Token).GetValue()
	that.Slock.Lock()
	defer that.Slock.Unlock()
	that.Data = make(map[string]interface{})
	json.Unmarshal([]byte(datatype.Type2str(data)), &that.Data)
	return that
}

func (that *WebAccessToken) DelToken() {
	that.Redis.SetKey(that.Token).DelKey()
}

func (that *WebAccessToken) GetTokenValue(key string) interface{} {
	that.GetTokenData()
	that.Slock.Lock()
	defer that.Slock.Unlock()
	if common.Has_map_index(key, that.Data) {
		return that.Data[key]
	} else {
		return nil
	}
}

func (that *WebAccessToken) UpdateToken() {
	if that.Redis.SetKey(that.Token).HasKey() {
		that.Redis.SetKey(that.Token).UpdateTTL(that.TimeOut)
	}
}

func (that *WebAccessToken) SetData(data map[string]interface{}) {
	that.Slock.Lock()
	defer that.Slock.Unlock()
	that.Data = make(map[string]interface{})
	for key, val := range data {
		that.Data[key] = val
	}
	j_data, _ := json.Marshal(that.Data)
	that.Redis.SetKey(that.Token).SetExValue(string(j_data), that.TimeOut)
}

func (that *WebAccessToken) SetList(data map[string]interface{}) {
	that.Slock.Lock()
	defer that.Slock.Unlock()
	for key, val := range data {
		that.Data[key] = val
	}
	j_data, _ := json.Marshal(that.Data)
	that.Redis.SetKey(that.Token).SetList(string(j_data))
}

func (that *WebAccessToken) StartNewToken() *WebAccessToken {
	uu_id := common.Get_UUID()
	that.Token = uu_id
	return that
}
