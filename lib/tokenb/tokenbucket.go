package tokenb

import (
	"time"
	"sync"
)


type TokenBucket struct {
	rate   int           // tokens per second (RPS)
	burst  int           // max tokens
	tokens chan struct{}
	stop   chan struct{}
}


type RampingBucket struct {
	rate      RateFunc
	start     time.Time
	nextToken time.Time
	burst     int

	tokens int
	mu     sync.Mutex
}


func NewTokenBucket(rate, burst int) *TokenBucket {
	tb := &TokenBucket{
		rate:   rate,
		burst: burst,
		tokens: make(chan struct{}, burst),
		stop:   make(chan struct{}),
	}
	go tb.run()
	return tb
}


func (tb *TokenBucket) run() {
	interval := time.Second / time.Duration(tb.rate)
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case tb.tokens <- struct{}{}:
			default:
				// bucket full → token dropped
			}
		case <-tb.stop:
			return
		}
	}
}

func (tb *TokenBucket) Take() {
	<-tb.tokens
}


func NewRampingBucket(rate RateFunc, burst int) *RampingBucket {
	now := time.Now()

	return &RampingBucket{
		rate:      rate,
		start:     now,
		nextToken: now,
		burst:     burst,
	}
}

func (b *RampingBucket) Take() {
	for {
		b.mu.Lock()

		now := time.Now()

		// 1. Пополняем токены
		elapsed := now.Sub(b.start)
		rps := b.rate(elapsed)

		if rps <= 0 {
			b.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			continue
		}

		interval := time.Duration(float64(time.Second) / rps)

		for !b.nextToken.After(now) && b.tokens < b.burst {
			b.tokens++
			b.nextToken = b.nextToken.Add(interval)
		}

		// 2. Если токен есть — забираем
		if b.tokens > 0 {
			b.tokens--
			b.mu.Unlock()
			return
		}

		// 3. Ждём до следующего токена
		sleep := b.nextToken.Sub(now)
		b.mu.Unlock()

		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}


