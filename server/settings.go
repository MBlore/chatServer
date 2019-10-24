package server

import "time"

const (
	readTimeoutDuration  = 60 * time.Second
	writeTimeoutDuration = 30 * time.Second
)
