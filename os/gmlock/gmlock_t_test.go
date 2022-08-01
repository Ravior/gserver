package gmlock

import (
	"log"
	"testing"
)

func Test_Lock(t *testing.T) {
	lockKey := "lock"
	Lock(lockKey)
	if TryLock(lockKey) == true {
		t.Fail()
	}
	Unlock(lockKey)
	if TryLock(lockKey) == false {
		t.Fail()
	}
	Unlock(lockKey)
	if TryLockFunc(lockKey, func() {
		log.Println("加锁成功")
	}) == false {
		t.Fail()
	}
	Unlock(lockKey)
}

func Test_RLock(t *testing.T) {
	lockKey := "lock"
	RLock(lockKey)
	RUnlock(lockKey)
	Lock(lockKey)
	if TryRLock(lockKey) == true {
		t.Fail()
	}
	Unlock(lockKey)
}
