package heartbeat

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/cenk/backoff"
)

/* start of a check interface...
type Check interface {
	Beat
	Probe
}
*/

// Check Defines a standard heartbeat Check
type Check struct {
	URL      *url.URL
	Retries  int
	Timeout  time.Duration
	Interval time.Duration
	OK       bool
	sync.Mutex
}

// MaxRetryBackOff Implements an implementation of backoff that has a max retry
type MaxRetryBackOff struct {
	Interval   time.Duration
	MaxRetries int
	tries      int
}

// Reset return to iniital state
func (b *MaxRetryBackOff) Reset() { b.tries = 0 }

// NextBackOff checks the max retries and returns stop if its reached its threshold
func (b *MaxRetryBackOff) NextBackOff() time.Duration {
	if b.tries >= b.MaxRetries {
		return backoff.Stop
	}

	b.tries++
	return b.Interval
}

// TODO: refactor this into a check interface  then make check for each scheme
func (c *Check) probeHTTP() error {
	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Get(c.URL.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("request returned %d expected code 200: %s", res.StatusCode, res.Status)
	}
	return nil
}

func (c *Check) check() error {
	conn, err := net.DialTimeout("tcp", c.URL.Host, c.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	// select on scheme and do some more poking
	switch {
	case c.URL.Scheme == "http":
		return c.probeHTTP()
	}

	return nil
}

// Beat - starts the heart beat checker.
func (c *Check) Beat() error {
	operation := func() error {
		return c.check()
	}

	ticker := time.NewTicker(c.Interval)
	for {
		select {
		case <-ticker.C:
			c.Lock()

			err := backoff.Retry(operation, &MaxRetryBackOff{Interval: c.Interval, MaxRetries: c.Retries})
			if err != nil {
				c.OK = false
				c.Unlock()
				return fmt.Errorf("\nFailed to connect after %d tries: %s\n", c.Retries, err)
			}
			log.Println("Check succeeded for ", c.URL.String())
			c.OK = true
			c.Unlock()
		}
	}
}
