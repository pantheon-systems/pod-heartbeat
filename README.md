[![Docker Repository on Quay](https://quay.io/repository/getpantheon/pod-heartbeat/status "Docker Repository on Quay")](https://quay.io/repository/getpantheon/pod-heartbeat)
[![Coverage Status](https://coveralls.io/repos/github/pantheon-systems/pod-heartbeat/badge.svg?branch=master)](https://coveralls.io/github/pantheon-systems/pod-heartbeat?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pantheon-systems/pod-heartbeat)](https://goreportcard.com/report/github.com/pantheon-systems/pod-heartbeat)

# pod heartbeat
Connect to services or exit allowing a pod to restart.

This is inteded to be used as a side car container in a kube pod to short-circuit a pods existence if it can't reach some service


## Usage

This command will connect to localhost on port 4000 with a 1 second timeout for the probe, it probe every 5 seconds, and will exit if the check fails 3 times 
```
  ./pod-heartbeat --connect http://localhost:4000 --timeout 1s --interval 5s --retries 3 
```

### using as a kube probe
The service can be invoked in 'server' mode. Where it reports 200 or 503 codes depending on health. you can couple this with kube liveness probes to restart pods.
```
 ./pod-heartbeat --connect http://localhost:4000 --server 
```

by default it will listen on *:9999 and the root handler will report status as a 200 or 503 which you can use in kube pods health check


## Help

Use -h / --help to see help
```
 ./pod-heartbeat --help                                                                                                                                                                                                                                                                                                                           130 ↵
Sometimes you want your kube pod to die if it can't get to something.
That could be another container in your pod that has deadlocked.

This program runs connects, and  if it can't connect or times out it will exit.
When ran inside a kube pod the container exit event will cause the pod to be destroyed.

Usage:
  pod-heartbeat [flags]

Flags:
  -f, --config-file string   Config file (default is $HOME/.pod-heartbeat.yaml)
  -c, --connect string       Connection URI, valid protocols are  tcp:// and http:// for now (default "tcp://127.0.0.1:4000")
  -i, --interval string      Interval for the Heartbeat action. (default "5s")
  -p, --port int             Port to listen on for the status server. (default 9999)
  -r, --retries int          How many times to retry before reporting failure. (default 3)
  -s, --server               Run the status server instead of exiting.
  -t, --timeout string       Timeout before considering the connection failed. Valid qualifiers: ns,ms,s,m,h,d (default "1s")
```
