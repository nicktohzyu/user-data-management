// Package pool based on https://github.com/go-home-iot/connection-pool
package pool

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
	"user-data-management/commons"
	"user-data-management/commons/logger"
)

// ErrTimeout represents a timeout error, for example you called Get and couldn't get
// a connection within the timeout period.
var ErrTimeout = errors.New("timeout")

// ConnectionPool provides the ability to pool connections
type ConnectionPool struct {
	Config  Config
	pool    chan *net.Conn
	closed  bool
	address string
}

// NewPool creates a new ConnectionPool.  The pool which is returned will still need to
// have Init() called in it before it can be used
func NewPool(config Config, address string) *ConnectionPool {
	p := &ConnectionPool{
		Config:  config,
		pool:    make(chan *net.Conn, config.NumConnections),
		closed:  false,
		address: address,
	}
	p.Init()
	return p
}

// Init should be called before using the pool, the call is non blocking, but you
// can wait on the returned channel if you want to know when all of the underlying
// connections have been created and are ready to use
func (p *ConnectionPool) Init() {
	loggerLevel := logger.Level
	logger.Level = logger.INFO
	logger.Info("Initializing connection pool")
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(p.Config.NumConnections)

	for i := 0; i < p.Config.NumConnections; i++ {
		p.retryNewConnection(&wg)
	}

	// Return the channel to let the caller know when init has completed
	wg.Wait()
	endTime := time.Now()
	logger.Info("Connection pool initialized. Time taken: ", endTime.Sub(startTime))
	defer func() { logger.Level = loggerLevel }()
}

// Close closes all of the underlying connections, this is non blocking but you can
// wait on the returned channel if you need to know all the connections have closed
func (p *ConnectionPool) Close() chan bool {
	done := make(chan bool)
	go func() {
		for len(p.pool) > 0 {
			c := <-p.pool
			(*c).Close()
		}
		done <- true
	}()
	return done
}

// Get is a blocking function that waits to get an available connection.  If after the
// timeout duration a connection could not be fetched, the function returns with ErrTimeout.
// The flush parameter if set to true will read all of the outstanding data from the
// connection before returning it to the caller. Note there is a possible 100ms delay for this
// function to return if you set flush==true while the pool tries to read any existing content
// from the connection
func (p *ConnectionPool) Get() (*net.Conn, error) {
	startTime := time.Now()
	expire := time.Now().Add(p.Config.WaitDuration)
	var (
		c   *net.Conn = nil
		err error     = nil
	)
	select {
	case conn := <-p.pool:
		//if flush {
		//	// Read all the contents from the buffer, if there is any, then
		//	// reset the read deadline to infinity
		//	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		//	_, _ = ioutil.ReadAll(conn)
		//	conn.SetReadDeadline(time.Time{})
		//}
		c, err = conn, nil
		break
	case <-time.After(expire.Sub(time.Now())):
		c, err = nil, ErrTimeout
	}
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ConnPoolComponent, commons.GetConnLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return c, err
}

// Release returns the connection back to the pool. err is any error that was returned
// by the connection while it was being used, if there was an error the pool will then
// throw this connection away and create a new one
func (p *ConnectionPool) Release(c *net.Conn) {
	startTime := time.Now()
	if c == nil {
		logger.Error("Warning: releasing nil connection")
		return
	}
	p.pool <- c

	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ConnPoolComponent, commons.FreeConnLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
}

// NewConnection returns an initialized Connection instance
func NewConnection(address string, p *ConnectionPool) (*net.Conn, error) {
	logger.Debug("Pool creating connection to: ", address)
	c, err := net.DialTimeout("tcp", address, p.Config.WaitDuration)
	if err != nil {
		logger.Error("Connection pool new connection error: ", err)
		return nil, err
	}
	return &c, nil
}

func (p *ConnectionPool) retryNewConnection(wg *sync.WaitGroup) {
	// Just keeps trying to open a new connection until it succeeds
	startTime := time.Now()
	func() {
		for !p.closed {
			//TODO: set max number of attempts
			//logger.Info("Pool attempting to create new connection")
			c, err := NewConnection(p.address, p)
			if err == nil {
				//logger.Info("Connection created")
				p.pool <- c
				if wg != nil {
					wg.Done()
				}
				return
			}
			logger.Info("Error creating new connection:", err)
			// Wait for a small time then retry
			time.Sleep(p.Config.RetryDuration)
		}
	}()
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ConnPoolComponent, commons.InitConnLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
}

func (p *ConnectionPool) SendBytes(outBytes []byte) (*[]byte, error) {
	loggerLevel := logger.Level
	logger.Level = logger.ERROR
	defer func() { logger.Level = loggerLevel }()
	cp, err := p.Get()
	defer p.Release(cp) //TODO: use err to reinitialize problematic connection
	if err != nil {
		logger.Error("Error getting connection from pool:", err)
		return nil, err
	}
	c := *cp
	//TODO: clean output so that it does not contain \n
	err = c.SetDeadline(time.Now().Add(p.Config.ReadDuration))
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	logger.Debug("Sending string: ", string(outBytes))
	_, err = fmt.Fprintf(c, string(outBytes)+"\n")
	if err != nil {
		logger.Error("Error in sending")
		return nil, err
	}
	logger.Debug("String sent")
	endTime := time.Now()
	commons.Latency.WithLabelValues(
		commons.ConnPoolComponent, commons.SendMsgLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))

	startTime = time.Now()
	reader := bufio.NewReader(c)
	err = c.SetDeadline(time.Now().Add(p.Config.ReadDuration))
	if err != nil {
		return nil, err
	}
	logger.Debug("Reading response")
	response, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Error("Error reading response", err)
		return nil, err
	}

	logger.Debug("Conn pool received response: ", string(response))
	endTime = time.Now()
	commons.Latency.WithLabelValues(
		commons.ConnPoolComponent, commons.GetResponseLabel).
		Observe(float64(endTime.Sub(startTime).Milliseconds()))
	return &response, nil
}
