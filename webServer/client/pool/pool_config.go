package pool

import (
	"time"
)

// Config contains all of the configuration parameters for the connection pool
type Config struct {

	// NumConnections is the number of connections to open
	NumConnections int

	// RetryDuration specifies how long the pool will wait to try to create a new connection
	// if the previous new connection attempt failed
	RetryDuration time.Duration

	// WaitDuration specifies how long the pool's Get method will wait before returning ErrTimeout
	WaitDuration time.Duration

	// ReadDuration specifies how long the pool will wait to send and get a response
	ReadDuration time.Duration
}

const (
	NUM_CONNECTIONS = 1000
)
