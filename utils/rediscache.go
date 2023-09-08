package utils

import (
	_ "github.com/beego/beego/v2/client/cache/redis"
	"encoding/json"
	"github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"time"
	"context"
)
var RedisCache cache.Cache

func init(){
redishost,_:=	config.String("redis::redis_host")
redisport,_:=	config.String("redis::redis_port")
redisdb,_:=	config.String("redis::redis_db")
redispass,_:=	config.String("redis::redis_pass")
cacheRedisConn,_ := json.Marshal(map[string]string {
	"key" : "redisCache",
	"conn" : redishost+":"+redisport,
	"dbNum" :  redisdb,
	"password" : redispass,
 })
 log := logs.NewLogger()

 var err error
 RedisCache,err = cache.NewCache("redis", string(cacheRedisConn))
 if err != nil || RedisCache == nil {
	errMsg := "failed to init redis"
	log.Error(errMsg, err)
	// panic(errMsg)
}
}
///set cache key
func SetStr(key, value string, time time.Duration) (err error) {
	err = RedisCache.Put(context.Background(),key, value, time)
	if err != nil {
		log := logs.NewLogger()
		log.Error("set key:", key, ",value:", value, err)
	}
	return
}
///get cache key
func GetStr(key string) (value string,err error) {
	v,err := RedisCache.Get(context.Background(),key)
	if(err!=nil){
		return "",err
	}
	if(v==nil){//item not exist
		return "",nil
	}
	value = string(v.([]byte)) //这里的转换很重要，Get返回的是interface
	return value,nil
}
///delete cache key
func DelKey(key string) (err error) {
	err = RedisCache.Delete(context.Background(),key)
	if(err!=nil){
		return err
	}
	return nil
}

