package limiter

import "fmt"

type RequestLimiter struct {
	maxConcurrentRequest int
	tokenPool            chan int
}

func NewRequestLimiter(maxCR int) *RequestLimiter {
	return &RequestLimiter{maxConcurrentRequest: maxCR, tokenPool: make(chan int, maxCR)}
}

func (this *RequestLimiter) GetToken() bool {
	if len(this.tokenPool) >= this.maxConcurrentRequest {
		fmt.Println("here")
		return false
	} else {
		this.tokenPool <- 1
		return true
	}
}

func (this *RequestLimiter) ReleaseToken() {
	if len(this.tokenPool) == 0 {
		return
	}
	_ = <-this.tokenPool
	return
}
