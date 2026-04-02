package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/samthor/daikinac"
)

func main() {

	devices := []daikinac.Device{
		{Host: "192.168.3.145"}, // den
		{Host: "192.168.3.152"}, // living room
		{Host: "192.168.3.204"}, // bedroom
		{Host: "192.168.3.225"}, // loft
		{Host: "192.168.3.245", UUID: "f45aab28604811eca7c4737954d1686f"}, // office
		// with UUID takes >second per request; slow CPU doing fake SSL?
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	aggregateDevices(ctx, devices)
}

func aggregateDevices(ctx context.Context, d []daikinac.Device) {
	var wg sync.WaitGroup

	type status struct {
		Error error
		daikinac.Status
	}
	out := make([]status, len(d))

	for i, device := range d {
		wg.Add(1)

		go func() {
			defer wg.Done()

			now := time.Now()
			defer func() {
				log.Printf("addr %s took %v", device.Host, time.Since(now))
			}()

			s, err := device.FetchAll(ctx)
			if err != nil {
				out[i].Error = err
			}
			out[i].Status = s
		}()
	}

	wg.Wait()

	for _, s := range out {
		log.Printf("status=%+v", s)
	}
}
