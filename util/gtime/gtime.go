package gtime

import (
	"github.com/Ravior/gserver/util/gconv"
	"time"
)

const (
	SevenDaysSeconds int64 = 7 * 86400
	OneDaySeconds    int64 = 86400
)

// Now 当前时间
func Now() int64 {
	return time.Now().Unix()
}

func NowHour() int {
	return time.Now().Hour()
}

func NowMinute() int {
	return time.Now().Minute()
}

func NowWeek() int {
	nowWeek := int(time.Now().Weekday())
	if nowWeek == 0 {
		return 7
	}
	return nowWeek
}

// NowDate 当前日期 格式 Ymd
func NowDate() int32 {
	return gconv.Int32(time.Now().Format("20060102"))
}

// NowTimeStr 当前日期 格式 Y-m-d H:i:s
func NowTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// NowDateStr 当前日期 格式 Ymd
func NowDateStr() string {
	return time.Now().Format("20060102")
}

// YesterdayDateStr 昨天日期 格式 Ymd
func YesterdayDateStr() string {
	return time.Now().AddDate(0, 0, -1).Format("20060102")
}

// UnixMs 获取当前时间(毫秒）
func UnixMs() int64 {
	return time.Now().UnixNano() / 1e6
}

// UnixNano 获取当前时间(纳秒）
func UnixNano() int64 {
	return time.Now().UnixNano()
}

// IsToday 判断是不是今天
func IsToday(ts int64) bool {
	return time.Unix(ts, 0).Format("20060102") == time.Now().Format("20060102")
}

// IsYesterday 判断是不是昨天
func IsYesterday(ts int64) bool {
	return time.Unix(ts, 0).Format("20060102") == YesterdayDateStr()
}

// StartOfDay 当天开始时间戳
func StartOfDay(ts int64) int64 {
	formatTimeStr := time.Unix(ts, 0).Format("2006-01-02")
	formatTime, _ := time.ParseInLocation("2006-01-02", formatTimeStr, time.Local)
	return formatTime.Unix()
}

// NowAddDays 当前时间+多少天
func NowAddDays(day int64) int64 {
	return Now() + day*OneDaySeconds
}
