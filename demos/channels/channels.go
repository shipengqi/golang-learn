package channels

const (
	max = 50000000
	block = 500
	bufsize = 100
)

func multiSend() {
	done := make(chan struct{})
	c := make(chan int, bufsize)

	go func() {
		count := 0
		for x := range c {
			count += x
		}
		close(done)
	}()

	for i := 0; i < max; i ++ {
		c <- i
	}
	close(c)
	<-done
}

func blockSend() {
	done := make(chan struct{})
	c := make(chan [block]int, bufsize)

	go func() {
		count := 0
		for a := range c {
			for _, x := range a {
				count += x
			}
		}
		close(done)
	}()

	for i := 0; i < max; i += block {
		var b [block]int
		for n :=  0; n < block; n ++ {
			b[n] = i + n
			if i + n == max - 1 {
				break
			}
		}
		c <- b
	}
	close(c)
	<-done
}
