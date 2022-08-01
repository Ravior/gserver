package gttllock

import (
	"testing"
	"time"
)

func Test_TryLock(t *testing.T) {
	lockKey := "lock"
	TryLock(lockKey, 100)
	time.Sleep(50 * time.Millisecond)
	if TryLock(lockKey, 50) == true {
		t.Fail()
	}
	time.Sleep(100 * time.Millisecond)
	if TryLock(lockKey, 50) == false {
		t.Fail()
	}
}

func Test_IsLock(t *testing.T) {
	lockKey := "lock"
	TryLock(lockKey, 100)
	if IsLock(lockKey) == false {
		t.Fail()
	}
}

func Test_Unlock(t *testing.T) {
	lockKey := "lock"
	Lock(lockKey, 100)
	Unlock(lockKey)
	if IsLock(lockKey) == true {
		t.Fail()
	}
}
