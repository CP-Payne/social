package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}

// TODO: Try to implement the rate limiter in redis instead of using in-memory map
