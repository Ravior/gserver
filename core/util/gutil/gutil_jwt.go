package gutil

import (
	"fmt"
	"github.com/Ravior/gserver/core/util/gconv"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Jwt struct {
	Field string
	Key   []byte
	Ttl   int64
}

func (j *Jwt) Encode(fieldValue interface{}) string {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		j.Field: fieldValue,
		"exp":   time.Now().Unix() + j.Ttl,
	})
	token, err := at.SignedString(j.Key)
	if err != nil {
		fmt.Println("JWT Encode Error:", err)
		return ""
	}
	return token
}

func (j *Jwt) Decode(auth string) (fieldValue interface{}, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = claim.Claims.(jwt.MapClaims)[j.Field]
	return
}

func (j *Jwt) DecodeUint(auth string) (fieldValue uint, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Uint(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeUint8(auth string) (fieldValue uint8, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Uint8(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeUint32(auth string) (fieldValue uint32, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Uint32(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeUint64(auth string) (fieldValue uint64, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Uint64(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeInt(auth string) (fieldValue int, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Int(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeInt8(auth string) (fieldValue int8, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Int8(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeInt32(auth string) (fieldValue int32, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Int32(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeInt64(auth string) (fieldValue int64, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.Int64(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}

func (j *Jwt) DecodeString(auth string) (fieldValue string, err error) {
	claim, err := jwt.Parse(auth, func(token *jwt.Token) (i interface{}, err error) {
		return j.Key, nil
	})
	if err != nil {
		// 解析失败
		fmt.Println("JWT DecodeUint32 Error:", err)
		return
	}

	fieldValue = gconv.String(claim.Claims.(jwt.MapClaims)[j.Field])
	return
}
