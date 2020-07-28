package dq

import (
	"sync"

	"github.com/beanstalkd/go-beanstalk"
)

type connection struct {
	lock     sync.RWMutex
	endpoint string
	tube     string
	conn     *beanstalk.Conn
}

func newConnection(endpint, tube string) *connection {
	return &connection{
		endpoint: endpint,
		tube:     tube,
	}
}

func (c *connection) Close() error {
	c.lock.Lock()
	conn := c.conn
	c.conn = nil
	defer c.lock.Unlock()

	if conn != nil {
		return conn.Close()
	}

	return nil
}

func (c *connection) get() (*beanstalk.Conn, error) {
	c.lock.RLock()
	conn := c.conn
	c.lock.RUnlock()
	if conn != nil {
		return conn, nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	var err error
	c.conn, err = beanstalk.Dial("tcp", c.endpoint)
	if err != nil {
		return nil, err
	}

	c.conn.Tube.Name = c.tube
	return c.conn, err
}

func (c *connection) reset() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}
