

## Design
There are two layers, a frontend layer that receives outside http and https requests and uses consistent hashing to route to a set of backend servers and a backend layer which saves 

```bash
./dpc_frontend.go --shmId my-app-name --hosts-file /tmp/dpc_hosts
[cachehit https://localhost:8000/foo]
[cachemiss https://localhost:8000/bar]
^c
1827 cache misses (38%)
340 cache hits (44%)
Hit Thread Stats   Avg      Stdev     Max 
  Latency   635.91us    0.89ms  12.92ms  
  Requests/sec   56.0k	 8.07k	 64k
Miss Thread Stats   Avg      Stdev     Max 
  Latency   1235.91us    0.89ms  12.92ms  
  Requests/sec   18.0k	 4.07k	 28k




./dpc_backend.go --port 87527 --threshold 8

dpc_hosts
localhost:87527
87528
127.0.0.1:87529
```
### MVP
Logging and statistics are extraneous. Dont have to configure thresholds, use static values first time around


what we want to measure:
the latency and thoroughput of the frontend, and the memory usage of the backend. The most basic thing we want to measure is how fast does it go when making requests via the proxy, vs when the proxy is missing. Test this by writing a lua script for wrk to randomly select from a large set of urls. Also test by using the debug-remote-port in a browser.
