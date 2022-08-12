package fungi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoneChanStream(t *testing.T) {
	ch := make(chan int)
	done := make(chan struct{})
	source := DoneChanStream(ch, done)

	var expected, sum int
	add := Loop(func(i int) { sum += i })
	go add(source)

	for i := 1; i < 1000; i++ {
		expected += i
		ch <- i
	}
	close(ch)
	done <- struct{}{}
	close(done)

	assert.Equal(t, expected, sum)
}

func TestChanStream(t *testing.T) {
	ch := make(chan int)
	source := ChanStream(ch)

	var expected, sum int
	add := Loop(func(i int) { sum += i })
	go add(source)

	for i := 1; i < 10; i++ {
		expected += i
		ch <- i
	}
	close(ch)

	assert.Equal(t, expected, sum)
}
