package jsonrpc

import (
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"golang.org/x/time/rate"
)

// RateLimit is the rate limit config
type RateLimit struct {
	rlm map[string]*rate.Limiter
	sync.RWMutex
}

var rateLimit = &RateLimit{}

// InitRateLimit initializes the rate limit config
func InitRateLimit(rlc RateLimitConfig) {
	setRateLimit(rlc)
}

// setRateLimit sets the rate limit config
func setRateLimit(rlc RateLimitConfig) {
	rateLimit.Lock()
	defer rateLimit.Unlock()
	rateLimit.rlm = updateRateLimit(rlc)
}

// updateRateLimit updates the rate limit config
func updateRateLimit(rateLimit RateLimitConfig) map[string]*rate.Limiter {
	log.Infof("rate limit config updated, config: %+v", rateLimit)
	if rateLimit.Enabled {
		log.Infof("rate limit enabled, api: %v, count: %d, duration: %d", rateLimit.RateLimitApis, rateLimit.RateLimitCount, rateLimit.RateLimitDuration)
		rlm := make(map[string]*rate.Limiter)
		for _, api := range rateLimit.RateLimitApis {
			rlm[api] = rate.NewLimiter(rate.Limit(rateLimit.RateLimitCount), rateLimit.RateLimitDuration)
		}
		for _, api := range rateLimit.SpecialApis {
			log.Infof("special api rate limit enabled, api: %v, count: %d, duration: %d", api.Api, api.Count, api.Duration)
			rlm[api.Api] = rate.NewLimiter(rate.Limit(api.Count), api.Duration)
		}
		return rlm
	}
	return nil
}

// methodRateLimitAllow returns true if the method is allowed by the rate limit
func methodRateLimitAllow(method string) bool {
	rateLimit.RLock()
	rlm := rateLimit.rlm
	rateLimit.RUnlock()
	if rlm != nil && rlm[method] != nil && !rlm[method].Allow() {
		return false
	}
	return true
}

// RateLimitConfig has parameters to config the rate limit
type RateLimitConfig struct {

	// Enabled defines if the rate limit is enabled or disabled
	Enabled bool `mapstructure:"Enabled"`

	// RateLimitApis defines the apis that need to be rate limited
	RateLimitApis []string `mapstructure:"RateLimitApis"`

	// RateLimitBurst defines the maximum burst size of requests
	RateLimitCount int `mapstructure:"RateLimitCount"`

	// RateLimitDuration defines the time window for the rate limit
	RateLimitDuration int `mapstructure:"RateLimitDuration"`

	// SpecialApis defines the apis that need to be rate limited with special rate limit
	SpecialApis []RateLimitItem `mapstructure:"SpecialApis"`
}

// RateLimitItem defines the special rate limit for some apis
type RateLimitItem struct {

	// Api defines the api that need to be rate limited
	Api string `mapstructure:"Api"`

	// Count defines the maximum burst size of requests
	Count int `mapstructure:"Count"`

	// Duration defines the time window for the rate limit
	Duration int `mapstructure:"Duration"`
}
