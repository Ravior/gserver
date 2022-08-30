package guuid

import "testing"

func BenchmarkGetUuid(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		GetUUID()
	}
}
