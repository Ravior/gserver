package gttllock

var locker = Locker{}

// TryLock 尝试加锁，时间为毫秒
func TryLock(key string, ttl int64) bool {
	return locker.TryLock(key, ttl)
}

// IsLock 加锁，时间为毫秒
func IsLock(key string) bool {
	return locker.IsLock(key)
}

// Lock 加锁，时间为毫秒
func Lock(key string, ttl int64) {
	locker.Lock(key, ttl)
}

func Unlock(key string) {
	locker.Unlock(key)
}
