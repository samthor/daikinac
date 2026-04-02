package daikinac

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// DoSetControl sets the ControlInfo on this Device.
func (d *Device) DoSetControl(c context.Context, in ControlInfo) (err error) {
	return d.Do(c, "/aircon/set_control_info", in, nil)
}

// DoGetControl reads the ControlInfo from this Device.
func (d *Device) DoGetControl(c context.Context) (out ControlInfo, err error) {
	err = d.Do(c, "/aircon/get_control_info", nil, &out)
	return
}

// DoGetSensor reads the SensorInfo from this Device.
func (d *Device) DoGetSensor(c context.Context) (out SensorInfo, err error) {
	err = d.Do(c, "/aircon/get_sensor_info", nil, &out)
	return
}

// DoGetBasic gets BasicInfo from this Device.
func (d *Device) DoGetBasic(c context.Context) (out BasicInfo, err error) {
	err = d.Do(c, "/common/basic_info", nil, &out)
	return
}

// Status aggregates the different information types of an aircon.
type Status struct {
	BasicInfo   BasicInfo
	ControlInfo ControlInfo
	SensorInfo  SensorInfo
}

// FetchAll is a helper which reads most public information about an AC.
func (d *Device) FetchAll(ctx context.Context) (s Status, err error) {
	var eg errgroup.Group

	eg.Go(func() error {
		s.BasicInfo, err = d.DoGetBasic(ctx)
		return err
	})

	eg.Go(func() error {
		s.ControlInfo, err = d.DoGetControl(ctx)
		return err
	})

	eg.Go(func() error {
		s.SensorInfo, err = d.DoGetSensor(ctx)
		return err
	})

	err = eg.Wait()
	return
}
