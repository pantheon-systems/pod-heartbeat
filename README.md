[![Docker Repository on Quay](https://quay.io/repository/getpantheon/pod-heartbeat/status "Docker Repository on Quay")](https://quay.io/repository/getpantheon/pod-heartbeat)
[![Coverage Status](https://coveralls.io/repos/github/pantheon-systems/pod-heartbeat/badge.svg?branch=master)](https://coveralls.io/github/pantheon-systems/pod-heartbeat?branch=master)

# pod heartbeat
Connect to services or exit allowing a pod to restart.

This is inteded to be used as a side car container in a kube pod to short-circuit a pods existence if it can't reach some service


## Usage

This command will connect to localhost on port 4000 with a 1 second timeout for the probe, it probe every 5 seconds, and will exit if the check fails 3 times 
```
  ./pod-heartbeat --connect http://localhost:4000 --timeout 1s --interval 5s --retries 3 
```


