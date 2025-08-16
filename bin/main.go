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
		{Host: "192.168.3.146"},
		{Host: "192.168.3.152"}, // living
		{Host: "192.168.3.204"},
		{Host: "192.168.3.225"},
		{Host: "192.168.3.245", UUID: "f45aab28604811eca7c4737954d1686f"}, // with UUID takes >second per request; slow CPU doing fake SSL?
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	aggregateDevices(ctx, devices)

	// // officeDevice := daikinac.Device{Host: "192.168.3.245", UUID: "f45aab28604811eca7c4737954d1686f"}

	// livingRoomDevice := daikinac.Device{Host: "192.168.3.152"}

	// err := livingRoomDevice.Do(context.TODO(), "/aircon/set_control_info", &daikinac.ControlInfo{
	// 	Power:   daikinac.Off,
	// 	Mode:    daikinac.ModeHeat,
	// 	SetTemp: 23.0,
	// 	FanRate: daikinac.FanAuto,
	// 	FanDir:  daikinac.FanNone,
	// }, nil)
	// log.Printf("control failure: %v", err)

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
			s, err := device.FetchAll(ctx)
			if err != nil {
				out[i].Error = err
			}
			out[i].Status = s
		}()
	}

	wg.Wait()
	log.Printf("got statuses=%+v", out)
}
