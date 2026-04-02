Raw interface to Daikin AC units.
Provides direct API access (via HTTP GET/POST).

Internal method descriptions can be found [here](https://github.com/ael-code/daikin-control/wiki/API-System), or elsewhere online.

## Device Support

Older devices can be accessed just by IP address.

Newer devices might need you to register a UUID and include in the `Device` struct.
See "helpers/register.sh" to perform a (one-time) registration of the UUID.
These devices require SSL with a pinned cert, which can be 10x as slow vs. older devices.
However, requests made in parallel go faster.

## Sample

Run `logger.go` or `broadcast.go` for simple demos.

Or, to dial a specific URL and get its sensor info:

```go
package main

import (
  "log"
  "context"

  "github.com/samthor/daikinac"
)

func main() {
  device := &daikinac.Device{
    Host: "192.168.1.155",
  }

  status, err := device.FetchAll(context.Background())
  if err != nil {
    log.Fatalf("couldn't fetch info: %v", err)
  }
  log.Printf("got status: %+v", status)
}
```
