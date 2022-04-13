package internal

import "time"

type middlewareAccess struct {
	xRemaining  int64
	xLimit      int64
	xRetryAfter time.Duration
	isAccess    bool
}

func NewMiddlewareAccess(xRemaining, xLimit int64, xRetryAfter time.Duration, isAccess bool) *middlewareAccess {
	return &middlewareAccess{xRemaining: xRemaining,
		xLimit:      xLimit,
		xRetryAfter: xRetryAfter,
		isAccess:    isAccess}
}
