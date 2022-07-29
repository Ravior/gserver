package gdb

import (
	"fmt"
	"github.com/Ravior/gserver/core/util/gconfig"
	"github.com/Ravior/gserver/core/util/gconv"
	"github.com/garyburd/redigo/redis"
)

var (
	clients          = make(map[string]*RedisCli, 1) // 创建RedisCli
	DefaultRedisPool = "default"                     // 默认链接池
)

func InitRedis() {
	for pool, config := range gconfig.Global.Redis {
		clients[pool] = &RedisCli{
			conn: pool,
			pool: &redis.Pool{ //实例化一个连接池
				MaxIdle:     5,    // 最大空闲链接
				MaxActive:   1000, // 连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
				IdleTimeout: 300,  // 连接关闭时间 300秒 （300秒不使用自动关闭）
				Dial: func() (redis.Conn, error) { // 要连接的redis数据库
					return redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
				},
			},
			prefix: config.Prefix,
		}
	}
}

func GetRedis(pools ...string) *RedisCli {
	pool := DefaultRedisPool
	if len(pools) > 0 && pools[0] != "" {
		pool = pools[0]
	}
	if redisCli, ok := clients[pool]; ok {
		return redisCli
	}
	return nil
}

type RedisCli struct {
	conn   string
	pool   *redis.Pool
	prefix string
}

func (r *RedisCli) getConn() redis.Conn {
	return r.pool.Get()
}

func (r *RedisCli) GetKeyWithPrefix(key string) string {
	return r.prefix + key
}

func (r *RedisCli) Exists(key string) (bool, error) {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(key)
	res, err := conn.Do("EXISTS", key)
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

func (r *RedisCli) Get(key string) (interface{}, error) {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(key)

	v, err := conn.Do("GET", key)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (r *RedisCli) GetString(key string) (string, error) {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(key)

	result, err := conn.Do("GET", key)
	if err != nil {
		return gconv.String(result), err
	}
	return "", nil
}

func (r *RedisCli) Set(key string, value interface{}, timeout ...uint64) error {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(key)

	if len(timeout) > 0 {
		_, err := conn.Do("SET", key, value, timeout)
		return err
	}
	_, err := conn.Do("SET", key, value)
	return err
}

func (r *RedisCli) Del(key string) (bool, error) {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(key)
	res, err := conn.Do("DEL", key)
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

func (r *RedisCli) Lock(key string, ttl int) bool {
	conn := r.getConn()
	defer conn.Close()

	key = r.GetKeyWithPrefix(fmt.Sprintf("lock:%s", key))
	t, _ := conn.Do("SET", key, 1, "NX", "EX", ttl)
	return t != nil
}

func (r *RedisCli) UnLock(key string) {
	key = r.GetKeyWithPrefix(fmt.Sprintf("lock:%s", key))
	_, _ = r.Del(key)
}
