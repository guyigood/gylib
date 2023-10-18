package redispack

import (
	"github.com/garyburd/redigo/redis"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	"time"
)

var Redis_data map[string]string
var RePool *redis.Pool

func Redisinit() {
	//RePool=nil
	Redis_data = make(map[string]string)
	Redis_data = common.Getini("conf/app.ini", "redis", map[string]string{"redis_host": "127.0.0.1", "redis_port": "6379", "redis_auth": "", "redis_db": "0", "redis_perfix": "", "redis_minpool": "5", "redis_maxpool": "20", "timeout": "60"})
	RePool = Set_redis_pool()
}

func Get_redis_pool() *redis.Pool {
	//if(RePool.ActiveCount()<=0
	/*if(RePool.ActiveCount()<=0){
		RePool=Set_redis_pool()
	}*/
	if RePool == nil {
		Redisinit()
	}
	return RePool
	//return Set_redis_pool()
}

func Set_redis_pool() *redis.Pool {
	data := Redis_data
	timeout := datatype.Str2Int(Redis_data["timeout"])
	return &redis.Pool{
		MaxIdle:     datatype.Str2Int(Redis_data["redis_minpool"]),
		MaxActive:   datatype.Str2Int(Redis_data["redis_maxpool"]),
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
	}
}
