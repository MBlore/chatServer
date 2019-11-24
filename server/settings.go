package server

import "time"

const (
	// The time in seconds to allow a socket read operation to complete.
	readTimeoutDuration = 60 * time.Second

	// The time in seconds to allow a socket write operation to complete.
	writeTimeoutDuration = 30 * time.Second

	// Maximum packet length size allowed. Anything higher is considered a DDOS and a client will be forcefully disconnected.
	maxPacketLength = 20000000 // 20mb - accounts for very large images.
)
