package gttllock

import (
	"github.com/Ravior/gserver/util/gtime"
	"sync"
)

type Locker struct {
	locks sync.Map
}

// TryLock 尝试加锁
func (l *Locker) TryLock(key string, ttl int64) bool {
	// 毫秒
	now := gtime.UnixMs()
	if lockEndTime, ok := l.locks.LoadOrStore(key, now+ttl); ok {
		if lockEndTime.(int64) >= now {
			return false
		} else {
			l.locks.Store(key, now+ttl)
		}
	}
	return true
}

// IsLock 是否已锁定
func (l *Locker) IsLock(key string) bool {
	// 毫秒
	now := gtime.UnixMs()
	if lockEndTime, ok := l.locks.Load(key); ok {
		if lockEndTime.(int64) >= now {
			return true
		}
	}
	return false
}

// Lock 加锁
func (l *Locker) Lock(key string, ttl int64) {
	// 毫秒
	now := gtime.UnixMs()
	// 锁定
	l.locks.Store(key, now+ttl)
}

// Unlock 解锁
func (l *Locker) Unlock(key string) {
	l.locks.Delete(key)
}
