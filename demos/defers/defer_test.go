package defers

import "testing"

func BenchmarkCall(b *testing.B)  {
	for i := 0; i < b.N; i ++ {
		call()
	}
}


func BenchmarkDeferCall(b *testing.B)  {
	for i := 0; i < b.N; i ++ {
		deferCall()
	}
}