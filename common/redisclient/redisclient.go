package redisclient

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"gylib/common"
	"gylib/common/datatype"
	"time"
)

var r_data map[string]string

type RedisClient struct {
	Key       string
	Field     string
	Re_prefix string
	RePool    *redis.Pool
}

func NewRedisCient() *RedisClient {
	this := new(RedisClient)
	this.RePool, this.Re_prefix = GetRedisPool()
	return this
}

func (that *RedisClient) CloseConnect() {
	that.RePool.Close()
}

func (that *RedisClient) DoCMD(cmdname string, arg ...interface{}) (interface{}, error) {
	client := that.RePool.Get()
	defer client.Close()
	replay, err := client.Do(cmdname, arg)
	return replay, err
}

func (that *RedisClient) Flushdb() {
	client := that.RePool.Get()
	//defer client.Close()
	client.Do("FLUSHDB")
	client.Close()
}

func (that *RedisClient) HasKey() bool {
	client := that.RePool.Get()
	//defer client.Close()
	hasok, err := client.Do("EXISTS", that.Key)
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

func (that *RedisClient) DelKey() bool {
	client := that.RePool.Get()
	//defer client.Close()
	_, err := client.Do("DEL", that.Key)
	client.Close()
	if err != nil {
		return false
	}

	return true
}

func (that *RedisClient) UpdateTTL(timeout int64) {
	client := that.RePool.Get()
	client.Do("Expire", that.Key, timeout)
	client.Close()
}

func (that *RedisClient) GetValue() interface{} {
	client := that.RePool.Get()
	//defer client.Close()
	raw, err := client.Do("GET", that.Key)
	client.Close()
	if err != nil {
		return nil
	}
	return raw
}

func (that *RedisClient) SetExValue(data interface{}, timeout int64) {
	client := that.RePool.Get()
	client.Do("SETEX", that.Key, timeout, data)
	client.Close()
}

func (that *RedisClient) SetValue(data interface{}) {
	client := that.RePool.Get()
	client.Do("SET", that.Key, data)
	client.Close()
}

func (that *RedisClient) GetKey() string {
	return that.Key
}

func (that *RedisClient) SetKey(key string) *RedisClient {
	if that.Re_prefix != "" {
		that.Key = that.Re_prefix + key
	} else {
		that.Key = key
	}
	//fmt.Println(this.Key,key)
	return that
}

func (that *RedisClient) SetList(data interface{}) {
	client := that.RePool.Get()
	//defer client.Close()
	//raw, _ := json.Marshal(&this.Data)
	_, err := client.Do("LPUSH", that.Key, data)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(this.Key,raw,client)
	client.Close()
}

func (that *RedisClient) Ltrim(start, stop int) {
	client := that.RePool.Get()
	//defer client.Close()
	client.Do("LTRIM", that.Key, start, stop)
	client.Close()
}

func (that *RedisClient) GetList() interface{} {
	client := that.RePool.Get()
	//defer client.Close()
	raw, err := client.Do("RPOP", that.Key)
	client.Close()
	if err != nil {
		return nil
	}
	return raw
}

func (that *RedisClient) HDel() {
	client := that.RePool.Get()
	defer client.Close()
	client.Do("HDEL", that.Key, that.Field)
}

func (that *RedisClient) HSetField(name string) *RedisClient {
	that.Field = name
	return that
}

func (that *RedisClient) HExists() bool {
	client := that.RePool.Get()
	//defer client.Close()
	hasok, err := client.Do("HEXISTS", that.Key, that.Field)
	//fmt.Println(hasok,that.Key,that.Field)
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

func (that *RedisClient) HSet(data interface{}) {
	client := that.RePool.Get()
	client.Do("HSET", that.Key, that.Field, data)
	client.Close()
}

func (that *RedisClient) HGet() interface{} {
	client := that.RePool.Get()
	//defer client.Close()
	raw, err := client.Do("HGET", that.Key, that.Field)
	client.Close()
	if err != nil {
		return nil
	}
	return raw
}

func (that *RedisClient) Hincby(n int) {
	client := that.RePool.Get()
	//defer client.Close()
	client.Do("HINCRBY", that.Key, that.Field, n)
	client.Close()

}

func (that *RedisClient) SAdd(data interface{}) { //添加集合元素
	client := that.RePool.Get()
	//defer client.Close()
	client.Do("SADD", that.Key, data)
	client.Close()
}

func (that *RedisClient) SRem(data interface{}) { //删除集合元素
	client := that.RePool.Get()
	//defer client.Close()
	client.Do("SREM", that.Key, data)
	client.Close()
}

func (that *RedisClient) SMembers() []interface{} { //获取集合列表
	client := that.RePool.Get()
	//defer client.Close()
	result, err := client.Do("SMEMBERS", that.Key)
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

func (that *RedisClient) SISMEMBER(data interface{}) bool { //判断成员元素(this.data)是否是集合this.key的成员
	client := that.RePool.Get()
	//defer client.Close()
	result, err := client.Do("SISMEMBER", that.Key, data)
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

func RedisInit() {
	r_data = make(map[string]string)
	r_data = common.Getini("conf/app.ini", "redis", map[string]string{"redis_host": "127.0.0.1", "redis_port": "6379", "redis_auth": "", "redis_db": "0", "redis_perfix": "", "redis_minpool": "5", "redis_maxpool": "20", "timeout": "60"})
}

func GetRedisPool() (*redis.Pool, string) {
	if r_data == nil {
		RedisInit()
	}
	data := r_data
	timeout := datatype.Str2Int(r_data["timeout"])
	return &redis.Pool{
		MaxIdle:     datatype.Str2Int(r_data["redis_minpool"]),
		MaxActive:   datatype.Str2Int(r_data["redis_maxpool"]),
		IdleTimeout: time.Duration(timeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", data["redis_host"]+":"+data["redis_port"])
			if err != nil {
				return nil, err
			}
			// 选择db
			if data["redis_auth"] != "" {
				c.Do("AUTH", data["redis_auth"])
			}
			c.Do("SELECT", data["redis_db"])

			return c, nil
		},
	}, data["redis_perfix"]
}
