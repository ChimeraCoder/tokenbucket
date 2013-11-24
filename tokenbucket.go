package tokenbucket

import (
	"time"
    "sync"
)

type bucket struct {
	capacity int64
	tokens   chan struct{}
    rate    time.Duration // Add a token to the bucket every 1/r units of time
    rateMutex sync.Mutex
}




func NewBucket(capacity int64, rate time.Duration) *bucket {

	//A bucket is simply a channel with a buffer representing the maximum size
	tokens := make(chan struct{}, capacity)

	b := &bucket{capacity, tokens, rate, sync.Mutex{}}

	//Set off a function that will continuously add tokens to the bucket
	go func(b *bucket) {
		for {
			b.tokens <- struct{}{}
            rate := b.GetRate()
			time.Sleep(rate)
		}
	}(b)

	return b
}


func (b *bucket) GetRate() time.Duration {
    b.rateMutex.Lock()
    tmp := b.rate
    b.rateMutex.Unlock()
    return tmp
}

func (b *bucket) SetRate(rate time.Duration) {
    b.rateMutex.Lock()
    b.rate = rate
    b.rateMutex.Unlock()
}



//AddTokens manually adds n tokens to the bucket
func (b *bucket) AddToken(n int64) {
}

func (b *bucket) withdrawTokens(n int64) error {
    for i := int64(0) ; i < n; i++{
        <-b.tokens
    }
    return nil
}



func (b *bucket) SpendToken(n int64) (<-chan error) {
    // Default to spending a single token
    if n < 0{
        n = 1
    }

    c := make(chan error)
    go func(b *bucket, n int64, c chan error){
        c <- b.withdrawTokens(n)
        close(c)
        return
    }(b, n, c)


    return c
}
