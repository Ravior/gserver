package gtimer

import (
	"log"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	timer := New()
	if timer == nil {
		t.Fail()
	}
}

func Test_Add(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.Add(500*time.Millisecond, func() {
		log.Println("schedule")
	})
	time.Sleep(5 * time.Second)
}

func Test_AddOnce(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.AddOnce(500*time.Millisecond, func() {
		log.Println("schedule1")
	})
	timer.AddOnce(500*time.Millisecond, func() {
		log.Println("schedule2")
	})
	time.Sleep(5 * time.Second)
}

func Test_AddTimes(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.AddTimes(500*time.Millisecond, 4, func() {
		log.Println("schedule")
	})
	time.Sleep(5 * time.Second)
}

func Test_DelayAdd(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.DelayAdd(3*time.Second, 500*time.Millisecond, func() {
		log.Println("schedule")
	})
	time.Sleep(5 * time.Second)
}

func Test_DelayAddOnce(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.DelayAddOnce(3*time.Second, 500*time.Millisecond, func() {
		log.Println("schedule")
	})
	time.Sleep(5 * time.Second)
}

func Test_After(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.After(3*time.Second, func() {
		log.Println("schedule")
	})
	time.Sleep(5 * time.Second)
}

func Test_Clear(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.After(3*time.Second, func() {
		log.Println("schedule")
	})
	timer.Clear()
	timer.After(3*time.Second, func() {
		log.Println("schedule2")
	})
	time.Sleep(5 * time.Second)
}

func Test_Close(t *testing.T) {
	timer := New()
	log.Println("start")
	timer.After(3*time.Second, func() {
		log.Println("schedule")
	})
	timer.Close()
	timer.After(3*time.Second, func() {
		log.Println("schedule2")
	})
	time.Sleep(5 * time.Second)
}
