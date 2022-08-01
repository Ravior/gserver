package gtime

import (
	"fmt"
	"testing"
	"time"
)

func Test_NowHour(t *testing.T) {
	fmt.Println(NowHour())
}

func Test_NowMinute(t *testing.T) {
	fmt.Println(NowMinute())
}

func Test_NowWeek(t *testing.T) {
	fmt.Println(NowWeek())
}

func Test_NowDate(t *testing.T) {
	fmt.Println(NowDate())
}

func Test_NowDateStr(t *testing.T) {
	fmt.Println(NowDateStr())
}

func Test_NowTimeStr(t *testing.T) {
	fmt.Println(NowTimeStr())
}

func Test_YesterdayDateStr(t *testing.T) {
	fmt.Println(YesterdayDateStr())
}

func Test_UnixMs(t *testing.T) {
	fmt.Println(UnixMs())
}

func Test_UnixNano(t *testing.T) {
	fmt.Println(UnixNano())
}

func Test_IsToday(t *testing.T) {
	if IsToday(time.Now().Unix()) == false {
		t.Fail()
	}

	if IsToday(time.Now().Unix()-86400) == true {
		t.Fail()
	}
}

func Test_IsYesterday(t *testing.T) {
	if IsYesterday(time.Now().Unix()) == true {
		t.Fail()
	}

	if IsYesterday(time.Now().Unix()-86400) == false {
		t.Fail()
	}
}

func Test_StartOfDay(t *testing.T) {
	fmt.Println(StartOfDay(time.Now().Unix()))
	if IsToday(StartOfDay(time.Now().Unix())) == false {
		t.Fail()
	}

	if IsYesterday(StartOfDay(time.Now().Unix())-86400) == false {
		t.Fail()
	}
}

func Test_NowAddDays(t *testing.T) {
	if IsYesterday(NowAddDays(-1)) == false {
		t.Fail()
	}
}
