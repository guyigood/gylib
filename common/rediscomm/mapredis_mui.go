package rediscomm

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/redispack"
)

type RedisComm_mui struct {
	Key         string
	RdPack      *redispack.RedisPack
	RePool      *redis.Pool
	Field       string
	Common_exec string
	Timeout     int
	Re_prefix   string
	Data        interface{}
}

func NewRedisComm_mui(action string) *RedisComm_mui {
	this := new(RedisComm_mui)
	this.RdPack = redispack.Redis_mui_init(action)
	this.RePool = this.RdPack.RePool
	this.Re_prefix = this.RdPack.Rp_data["redis_perfix"]
	this.Timeout = 3600
	this.Common_exec = "SET"
	return this
}

func (this *RedisComm_mui) Flushdb() {
	client := this.RePool.Get()
	client.Do("FLUSHDB")
	client.Close()
}

func (this *RedisComm_mui) HasKey() bool {
	client := this.RePool.Get()
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

func (this *RedisComm_mui) DelKey() bool {
	client := this.RePool.Get()
	_, err := client.Do("DEL", this.Key)
	client.Close()
	if err != nil {
		return false
	}

	return true
}

func (this *RedisComm_mui) GetRawValue() interface{} {
	client := this.RePool.Get()
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

func (this *RedisComm_mui) SetRawValue() {
	client := this.RePool.Get()
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, this.Data)
	} else {
		client.Do("SET", this.Key, this.Data)
	}
	client.Close()
}

func (this *RedisComm_mui) Getkey() string {
	return this.Key
}

func (this *RedisComm_mui) SetKey(key string) *RedisComm_mui {
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

func (this *RedisComm_mui) GetHmapLen() int {
	client := this.RePool.Get()
	raw, err := client.Do("HLEN", this.Key)
	client.Close()
	if err != nil {
		return 0
	} else {
		return datatype.Type2int(raw)
	}

}

func (this *RedisComm_mui) SetFiled(key string) *RedisComm_mui {
	this.Field = key
	return this
}

func (this *RedisComm_mui) SetExec(key string) *RedisComm_mui {
	this.Common_exec = key
	return this
}

func (this *RedisComm_mui) SetTime(timect int) *RedisComm_mui {
	this.Timeout = timect
	return this
}

func (this *RedisComm_mui) SetData(data interface{}) *RedisComm_mui {
	this.Data = data
	return this
}

func (this *RedisComm_mui) Get_value() interface{} {
	client := this.RePool.Get()
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

func (this *RedisComm_mui) Set_value() {
	client := this.RePool.Get()
	//defer client.Close()
	raw, _ := json.Marshal(&this.Data)
	if this.Common_exec == "SETEX" {
		client.Do("SETEX", this.Key, this.Timeout, raw)
	} else {
		client.Do("SET", this.Key, raw)
	}
	client.Close()
}

func (this *RedisComm_mui) SetList() {
	client := this.RePool.Get()
	//defer client.Close()
	//client.Close()
	raw, _ := json.Marshal(&this.Data)
	client.Do("LPUSH", this.Key, raw)
	client.Close()
}

func (this *RedisComm_mui) GetList() interface{} {
	client := this.RePool.Get()
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

func (this *RedisComm_mui) Hset_map() {
	client := this.RePool.Get()
	//defer client.Close()
	//defer client.Close()
	raw, _ := json.Marshal(&this.Data)
	client.Do("HSET", this.Key, this.Field, raw)
	client.Close()
}

func (this *RedisComm_mui) Hset_raw() {
	client := this.RePool.Get()
	//defer client.Close()
	raw := this.Data
	client.Do("HSET", this.Key, this.Field, raw)
	client.Close()
}

func (this *RedisComm_mui) Hdel_map() {
	client := this.RePool.Get()
	//defer client.Close()
	client.Do("HDEL", this.Key, this.Field)
	client.Close()
}

func (this *RedisComm_mui) HEXISTS() bool {
	client := this.RePool.Get()
	//defer client.Close()
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

func (this *RedisComm_mui) Hget_map() interface{} {
	client := this.RePool.Get()
	//defer client.Close()
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

func (this *RedisComm_mui) Hget_raw() interface{} {
	client := this.RePool.Get()
	//defer client.Close()
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

func (this *RedisComm_mui) Push(channel_name, message string) int { //发布者
	client := this.RePool.Get()
	//defer client.Close()
	//defer client.Close()
	raw, _ := client.Do("PUBLISH", channel_name, message)
	client.Close()
	return datatype.Type2int(raw)
}

func (this *RedisComm_mui) Hincby(n int) {
	client := this.RePool.Get()
	//defer client.Close()
	client.Do("HINCRBY", this.Key, this.Field, n)
	client.Close()

}
