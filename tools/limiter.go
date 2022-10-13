package tools

import (
	"sync"
	"time"
)

// Limiter 限流器
// 目前只做简单实现，之后可以扩展，当一定时间内多此触发限流将会直接抛弃请求
type Limiter struct {
	userMap  map[int64]time.Time //表示该事件之前的请求不会处理，会回复过于频繁
	duration time.Duration
	mu       sync.Mutex
}

func NewLimiter(d time.Duration) *Limiter {
	return &Limiter{
		userMap:  make(map[int64]time.Time),
		duration: d,
		mu:       sync.Mutex{},
	}
}

// Filter 表示是否同意放行请求
func (l *Limiter) Filter(qqNumber int64) bool {
	l.mu.Lock()
	if t, ok := l.userMap[qqNumber]; !ok {
		l.userMap[qqNumber] = time.Now().Add(l.duration)
		l.mu.Unlock()
		return true
	} else {
		l.userMap[qqNumber] = time.Now().Add(l.duration)
		l.mu.Unlock()
		//如果当前时间小于下一次之前的话，说明用户过于频繁
		if time.Now().Before(t) {
			return false
		}
		return true
	}
}
