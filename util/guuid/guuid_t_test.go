package guuid

import (
	"fmt"
	"sync"
	"testing"
)

func Test_GetUUID(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(GetUUID())
			wg.Done()
		}()
	}
	wg.Wait()
}
