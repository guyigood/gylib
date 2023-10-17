package redispack

import (
	"github.com/garyburd/redigo/redis"
	"gylib/common"
	"gylib/common/datatype"
	"time"
)

type RedisPack struct {
	Rp_data map[string]string
	RePool  *redis.Pool
}

func Redis_mui_init(action string) *RedisPack {
	this := new(RedisPack)
	this.Rp_data = make(map[string]string)
	ac_name := action
	if ac_name == "" {
		ac_name = "redis"
	}
	this.Rp_data = common.Getini("conf/app.ini", ac_name, map[string]string{"redis_host": "127.0.0.1", "redis_port": "6379", "redis_auth": "", "redis_db": "0", "redis_perfix": "", "redis_minpool": "5", "redis_maxpool": "20"})
	this.RePool = this.Set_redis_pool()
	return this
}

func (this *RedisPack) Get_redis_pool() *redis.Pool {
	/*if(this.RePool.ActiveCount()<=0){
		this.RePool=this.Set_redis_pool()
	}*/
	return this.RePool //this.Set_redis_pool()
}

func (this *RedisPack) Set_redis_pool() *redis.Pool {
	data := this.Rp_data
	return &redis.Pool{
		MaxIdle:     datatype.Str2Int(data["redis_minpool"]),
		MaxActive:   datatype.Str2Int(data["redis_maxpool"]),
		IdleTimeout: 180 * time.Second,
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
