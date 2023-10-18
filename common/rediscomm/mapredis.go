package rediscomm

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/redispack"
)

type RedisComm struct {
	Key         string
	Field       string
	Common_exec string
	Timeout     int
	Re_prefix   string
	Data        interface{}
	RedisPool   *redis.Pool
}

func NewRedisComm() *RedisComm {
	this := new(RedisComm)
	this.RedisPool = redispack.Get_redis_pool()
	this.Re_prefix = redispack.Redis_data["redis_perfix"]
	this.Timeout = 3600
	this.Common_exec = "SET"
	return this
}

func (this *RedisComm) CloseConnect() {
	this.RedisPool.Close()
}

func (this *RedisComm) Flushdb() {
	client := this.RedisPool.Get()
	//defer client.Close()
	client.Do("FLUSHDB")
	client.Close()
}

func (this *RedisComm) HasKey() bool {
	client := this.RedisPool.Get()
	//defer client.Close()
	hasok, err := client.Do("EXISTS", this.Key)
	client.Close()

	if err != nil {
		return false
	}
	if datatype.Type2int(hasok) == 0 {
		return false
	} else {
		return true
	}

}

func (this *RedisComm) DelKey() bool {
	client := this.RedisPool.Get()
	//defer client.Close()
	_, err := client.Do("DEL", this.Key)
	client.Close()
	if err != nil {
		return false
	}

	return true
}

func (this *RedisComm) GetRawValue() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("GET", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	return raw
}

func (this *RedisComm) SetRawValue() {
	client := this.RedisPool.Get()
	//defer client.Close()
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, this.Data)
	} else {
		client.Do("SET", this.Key, this.Data)
	}
	client.Close()

}

func (this *RedisComm) Getkey() string {
	return this.Key
}

func (this *RedisComm) SetKey(key string) *RedisComm {
	if this.Re_prefix != "" {
		strlen := len(this.Re_prefix)
		if len(key) < strlen {
			this.Key = this.Re_prefix + key
			return this
		}
		if key[:strlen] == this.Re_prefix {
			this.Key = key
		} else {
			this.Key = this.Re_prefix + key
		}
	} else {
		this.Key = key
	}
	//fmt.Println(this.Key,key)
	return this
}

func (this *RedisComm) GetHmapLen() int {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("HLEN", this.Key)
	client.Close()
	if err != nil {
		return 0
	} else {
		return datatype.Type2int(raw)
	}

}

func (this *RedisComm) SetFiled(key string) *RedisComm {
	this.Field = key
	return this
}

func (this *RedisComm) SetExec(key string) *RedisComm {
	this.Common_exec = key
	return this
}

func (this *RedisComm) SetTime(timect int) *RedisComm {
	this.Timeout = timect
	return this
}

func (this *RedisComm) SetData(data interface{}) *RedisComm {
	this.Data = data
	return this
}

func (this *RedisComm) GetMapValue() map[string]interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("GET", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	//var data interface{}
	data := make(map[string]interface{})
	err = json.Unmarshal(raw.([]byte), &data)
	if err != nil {
		return nil
	}
	return data
}

func (this *RedisComm) Get_value() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("GET", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	//var data interface{}
	err = json.Unmarshal(raw.([]byte), &this.Data)
	if err != nil {
		return nil
	}
	return this.Data
}

func (this *RedisComm) Set_value() {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, _ := json.Marshal(&this.Data)
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, raw)
	} else {
		client.Do("SET", this.Key, raw)

	}
	client.Close()
}

func (this *RedisComm) SetList() {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, _ := json.Marshal(&this.Data)
	_, err := client.Do("LPUSH", this.Key, raw)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(this.Key,raw,client)
	client.Close()
}

func (this *RedisComm) SetRawList() {
	client := this.RedisPool.Get()
	//defer client.Close()
	//raw, _ := json.Marshal(&this.Data)
	_, err := client.Do("LPUSH", this.Key, this.Data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(this.Key,raw,client)
	client.Close()
}

func (this *RedisComm) Ltrim(start, stop int) {
	client := this.RedisPool.Get()
	//defer client.Close()
	client.Do("LTRIM", this.Key, start, stop)
	client.Close()
}

func (this *RedisComm) GetList() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("RPOP", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	//var data interface{}
	err = json.Unmarshal(raw.([]byte), &this.Data)
	if err != nil {
		return nil
	}
	return this.Data
}

func (this *RedisComm) GetRawList() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("RPOP", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	return raw
}

func (this *RedisComm) Hset_map() {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, _ := json.Marshal(&this.Data)
	client.Do("HSET", this.Key, this.Field, raw)
	client.Close()
}

func (this *RedisComm) Hset_raw() {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw := this.Data
	client.Do("HSET", this.Key, this.Field, raw)
	client.Close()
}

func (this *RedisComm) Hdel_map() {
	client := this.RedisPool.Get()
	defer client.Close()
	client.Do("HDEL", this.Key, this.Field)
}

func (this *RedisComm) HEXISTS() bool {
	client := this.RedisPool.Get()
	//defer client.Close()
	hasok, err := client.Do("HEXISTS", this.Key, this.Field)
	client.Close()
	if err != nil {
		return false
	}
	if datatype.Type2int(hasok) == 0 {
		return false
	} else {
		return true
	}

}

func (this *RedisComm) GetDataMap() map[string]interface{} {
	if this.Data != nil {
		data, ok := this.Data.(map[string]interface{})
		if ok {
			return data
		}
	}
	return nil

}

func (this *RedisComm) Hget_map() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("HGET", this.Key, this.Field)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	//var data interface{}
	//fmt.Println("raw=",string(raw.([]byte)))
	err = json.Unmarshal(raw.([]byte), &this.Data)
	if err != nil {
		//fmt.Println(this.Key,this.Field,err)
		return nil
	}
	return this.Data
}

func (this *RedisComm) Hget_raw() interface{} {
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, err := client.Do("HGET", this.Key, this.Field)
	client.Close()
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}

	return raw
}

func (this *RedisComm) Push(channel_name, message string) int { //发布者
	client := this.RedisPool.Get()
	//defer client.Close()
	raw, _ := client.Do("PUBLISH", channel_name, message)
	client.Close()
	return datatype.Type2int(raw)
}

func (this *RedisComm) Hincby(n int) {
	client := this.RedisPool.Get()
	//defer client.Close()
	client.Do("HINCRBY", this.Key, this.Field, n)
	client.Close()

}

func (this RedisComm) SAdd() { //添加集合元素
	client := this.RedisPool.Get()
	//defer client.Close()
	client.Do("SADD", this.Key, this.Data)
	client.Close()
}

func (this *RedisComm) SRem() { //删除集合元素
	client := this.RedisPool.Get()
	//defer client.Close()
	client.Do("SREM", this.Key, this.Data)
	client.Close()
}

func (this *RedisComm) SMembers() []interface{} { //获取集合列表
	client := this.RedisPool.Get()
	//defer client.Close()
	result, err := client.Do("SMEMBERS", this.Key)
	client.Close()
	if err != nil {
		return nil
	}
	list, ok := result.([]interface{})
	if !ok {
		return nil
	}
	return list
}

func (this RedisComm) SISMEMBER() bool { //判断成员元素(this.data)是否是集合this.key的成员
	client := this.RedisPool.Get()
	//defer client.Close()
	result, err := client.Do("SISMEMBER", this.Key, this.Data)
	client.Close()
	if err != nil {
		return false
	}
	if datatype.Type2int(result) == 0 {
		return false
	} else {
		return true
	}
}
