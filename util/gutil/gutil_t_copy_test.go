package gutil

import (
	"reflect"
	"testing"
)

type person struct {
	id   int
	name string
	age  int
}

func Test_DeepCopy(t *testing.T) {
	p := &person{}
	p1 := DeepClone(p).(*person)
	if reflect.ValueOf(p).Pointer() == reflect.ValueOf(p1).Pointer() {
		t.Fail()
	}
}
