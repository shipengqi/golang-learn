package channels

import "testing"

func BenchmarkMultiSend(b *testing.B)  {
	for i := 0; i < b.N; i ++ {
		multiSend()
	}
}

func BenchmarkBlockSend(b *testing.B)  {
	for i := 0; i < b.N; i ++ {
		blockSend()
	}
}
