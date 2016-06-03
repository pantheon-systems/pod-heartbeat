package heartbeat_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/pantheon-systems/pod-heartbeat/pkg/heartbeat"
)

func TestMaxRetryBackoff(t *testing.T) {
	b := heartbeat.MaxRetryBackOff{
		Interval:   1 * time.Second,
		MaxRetries: 3,
	}

	if b.NextBackOff() != time.Second {
		t.Error("invalid interval")
	}
}

func TestCheck(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/beat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK ")
	})

	testServer := httptest.NewServer(mux)

	u, err := url.Parse(testServer.URL)
	if err != nil {
		t.Fatal("couldn't parse test server URL")
	}

	c := heartbeat.Check{
		Retries:  3,
		Interval: 500 * time.Millisecond,
		Timeout:  100 * time.Millisecond,
		URL:      u,
	}

	err = c.Beat()
	if err == nil {
		t.Error("Expected failure to beat but instead it worked: ", testServer.URL)
	}

	u, err = url.Parse("http://not_valid_host:80")
	if err != nil {
		t.Fatal("wtf happned here:", err.Error())
	}

	err = c.Beat()
	if err == nil {
		t.Error("Expected failure to beat but instead it worked: ", testServer.URL)
	}

	u, err = url.Parse(fmt.Sprintf("%s/beat", testServer.URL))
	if err != nil {
		t.Fatal("couldn't parse test server URL")
	}

	c.URL = u
	// Fire this off in a goroutine then close the server. The routine shouldn't error until after
	ch := make(chan string)
	go func() {
		err := c.Beat()
		if err != nil {
			ch <- err.Error()
		}
	}()

	select {
	case <-time.After(2 * time.Second):
		testServer.Close()
	case res := <-ch:
		t.Error(res)
	}

}
