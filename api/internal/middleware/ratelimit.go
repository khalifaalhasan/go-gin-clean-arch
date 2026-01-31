package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

const (
	RateLimitMax      = 3
	RateLimitDuration = 1 * time.Minute
)

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		// Pipeline Redis untuk atomicity dan performance
		pipe := rdb.Pipeline()
		incr := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, RateLimitDuration)
		_, err := pipe.Exec(ctx)

		if err != nil {
			// Fail-open strategy: jika Redis down, log error tapi allow request
			// atau Fail-closed: return 500. Di sini kita pilih 500 demi keamanan.
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		count := incr.Val()
		remaining := RateLimitMax - count
		if remaining < 0 {
			remaining = 0
		}

		// Set Header standard
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", RateLimitMax))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", int64(RateLimitDuration.Seconds())))

		if count > RateLimitMax {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			return
		}

		c.Next()
	}
}