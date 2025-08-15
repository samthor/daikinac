package main

import (
	"context"
	"log"

	"github.com/samthor/daikinac"
)

func main() {

	// devices := []daikinac.Device{
	// 	{Host: "192.168.3.204"},
	// 	{Host: "192.168.3.225"},
	// 	{Host: "192.168.3.152"},
	// 	{Host: "192.168.3.146"},
	// 	{Host: "192.168.3.245", UUID: "f45aab28604811eca7c4737954d1686f"}, // with UUID takes >second per request; slow CPU doing fake SSL?
	// }

	// for _, d := range devices {
	// 	var err error

	// 	timeout := time.Second * 2
	// 	if d.UUID != "" {
	// 		timeout *= 3
	// 	}

	// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	// 	defer cancel()

	// 	var bi daikinac.BasicInfo
	// 	err = d.Fetch(ctx, "/common/basic_info", nil, &bi)
	// 	log.Printf("fetched info: %+v err=%v", bi, err)

	// 	var c daikinac.ControlInfo
	// 	err = d.Fetch(ctx, "/aircon/get_control_info", nil, &c)
	// 	log.Printf("fetched control: %+v err=%v", c, err)

	// 	var s daikinac.SensorInfo
	// 	err = d.Fetch(ctx, "/aircon/get_sensor_info", nil, &s)
	// 	log.Printf("fetched sensor: %+v err=%v", s, err)
	// }

	officeDevice := daikinac.Device{Host: "192.168.3.245", UUID: "f45aab28604811eca7c4737954d1686f"}

	err := officeDevice.Fetch(context.TODO(), "/aircon/set_control_info", &daikinac.ControlInfo{
		Power:   daikinac.Off,
		Mode:    4,
		SetTemp: 23.0,
		FanRate: 1,
		FanDir:  0,
	}, nil)
	log.Printf("control failure: %v", err)

}
