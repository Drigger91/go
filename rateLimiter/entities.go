package main

type RateLimiter interface {
	Allow(id string)
}
