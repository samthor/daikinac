package daikinac

import (
	"context"
)

type Status struct {
	BasicInfo   BasicInfo
	ControlInfo ControlInfo
	SensorInfo  SensorInfo
}

func (d *Device) FetchAll(ctx context.Context) (s Status, err error) {
	err = d.Do(ctx, "/common/basic_info", nil, &s.BasicInfo)
	if err != nil {
		return
	}

	err = d.Do(ctx, "/aircon/get_control_info", nil, &s.ControlInfo)
	if err != nil {
		return
	}

	err = d.Do(ctx, "/aircon/get_sensor_info", nil, &s.SensorInfo)
	return
}
