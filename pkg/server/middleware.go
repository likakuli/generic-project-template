package server

import (
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

var (
	once sync.Once
	rl ratelimit.Limiter
)

func leakyBucketRateLimiter(qps int) gin.HandlerFunc {
	once.Do(func() {
		ratelimit.New(qps)
	})
	return func(c *gin.Context) {
		rl.Take()
	}
}
