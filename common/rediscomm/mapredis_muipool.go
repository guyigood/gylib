package rediscomm

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/redispack"
)

type RedisComm_pool struct {
	Key         string
	RdPack      *redispack.RedisPack
	clinet      redis.Conn
	is_con      bool
	Field       string
	Common_exec string
	Timeout     int
	Re_prefix   string
	Data        interface{}
}

func NewRedisComm_pool(action string) *RedisComm_pool {
	this := new(RedisComm_pool)
	this.RdPack = redispack.Redis_mui_init(action)
	this.clinet = this.RdPack.Get_redis_pool().Get()
	this.is_con = true
	this.Re_prefix = this.RdPack.Rp_data["redis_perfix"]
	this.Timeout = 3600
	this.Common_exec = "SET"
	return this
}

func (this *RedisComm_pool) GetPool() {
	if !this.is_con {
		this.clinet = this.RdPack.Get_redis_pool().Get()
	} else {
		if this.clinet == nil {
			this.clinet = this.RdPack.Get_redis_pool().Get()
		}
	}
}

func (this *RedisComm_pool) ClosePool() {
	this.clinet.Close()
	this.is_con = false
}

func (this *RedisComm_pool) Flushdb() {
	client := this.clinet
	client.Do("FLUSHDB")
}

func (this *RedisComm_pool) HasKey() bool {
	client := this.clinet
	hasok, err := client.Do("EXISTS", this.Key)
	if err != nil {
		return false
	}
	if datatype.Type2int(hasok) == 0 {
		return false
	} else {
		return true
	}

}

func (this *RedisComm_pool) DelKey() bool {
	client := this.clinet
	_, err := client.Do("DEL", this.Key)
	if err != nil {
		return false
	}

	return true
}

func (this *RedisComm_pool) GetRawValue() interface{} {
	client := this.clinet
	raw, err := client.Do("GET", this.Key)
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}
	return raw
}

func (this *RedisComm_pool) SetRawValue() {
	client := this.clinet
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, this.Data)
	} else {
		client.Do("SET", this.Key, this.Data)
	}

}

func (this *RedisComm_pool) Getkey() string {
	return this.Key
}

func (this *RedisComm_pool) SetKey(key string) *RedisComm_pool {
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

func (this *RedisComm_pool) GetHmapLen() int {
	client := this.clinet
	raw, err := client.Do("HLEN", this.Key)
	if err != nil {
		return 0
	} else {
		return datatype.Type2int(raw)
	}

}

func (this *RedisComm_pool) SetFiled(key string) *RedisComm_pool {
	this.Field = key
	return this
}

func (this *RedisComm_pool) SetExec(key string) *RedisComm_pool {
	this.Common_exec = key
	return this
}

func (this *RedisComm_pool) SetTime(timect int) *RedisComm_pool {
	this.Timeout = timect
	return this
}

func (this *RedisComm_pool) SetData(data interface{}) *RedisComm_pool {
	this.Data = data
	return this
}

func (this *RedisComm_pool) Get_value() interface{} {
	client := this.clinet
	raw, err := client.Do("GET", this.Key)
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

func (this *RedisComm_pool) Set_value() {
	client := this.clinet
	raw, _ := json.Marshal(&this.Data)
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, raw)
	} else {
		client.Do("SET", this.Key, raw)
	}
}

func (this *RedisComm_pool) SetList() {
	client := this.clinet
	raw, _ := json.Marshal(&this.Data)
	client.Do("LPUSH", this.Key, raw)
}

func (this *RedisComm_pool) GetList() interface{} {
	client := this.clinet
	raw, err := client.Do("RPOP", this.Key)
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

func (this *RedisComm_pool) Hset_map() {
	client := this.clinet
	raw, _ := json.Marshal(&this.Data)
	client.Do("HSET", this.Key, this.Field, raw)
}

func (this *RedisComm_pool) Hset_raw() {
	client := this.clinet
	raw := this.Data
	client.Do("HSET", this.Key, this.Field, raw)
}

func (this *RedisComm_pool) Hdel_map() {
	client := this.clinet
	client.Do("HDEL", this.Key, this.Field)
}

func (this *RedisComm_pool) HEXISTS() bool {
	client := this.clinet
	hasok, err := client.Do("HEXISTS", this.Key, this.Field)
	if err != nil {
		return false
	}
	if datatype.Type2int(hasok) == 0 {
		return false
	} else {
		return true
	}

}

func (this *RedisComm_pool) Hget_map() interface{} {
	client := this.clinet
	raw, err := client.Do("HGET", this.Key, this.Field)
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

func (this *RedisComm_pool) Hget_raw() interface{} {
	client := this.clinet
	raw, err := client.Do("HGET", this.Key, this.Field)
	if err != nil {
		return nil
	}
	if raw == nil {
		return nil
	}

	return raw
}

func (this *RedisComm_pool) Push(channel_name, message string) int { //发布者
	client := this.clinet
	raw, _ := client.Do("PUBLISH", channel_name, message)
	return datatype.Type2int(raw)
}

func (this *RedisComm_pool) Hincby(n int) {
	client := this.clinet
	client.Do("HINCRBY", this.Key, this.Field, n)

}
