# kill-process
Basic golang runner to kill all request process (hardcode in array for now).

# Architecture
Two go routines running in parallel:
- Lookup
- Kill
Lookup go routine would target the process to kill and send it to its channel.
Kill go routine would listen to the channel, fetch the process to kill and kill it via signal (safe kill).

# To run
```go run main.go```

# Build binary

``` go build ```
# Install binary

```go install```
