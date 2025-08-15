package daikinac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type SensorInfo struct {
	HomeTemp       float32 `json:"htemp"`
	HomeHumidity   float32 `json:"hhum"`
	OutsideTemp    float32 `json:"otemp"`
	CompressorFreq int     `json:"cmp"`
}

// ret=OK,pow=0,mode=4,adv=,stemp=19.0,shum=0,dt1=25.0,dt2=M,dt3=25.0,dt4=19.0,dt5=19.0,dt7=25.0,dh1=AUTO,dh2=50,dh3=0,dh4=0,dh5=0,dh7=AUTO,dhh=50,b_mode=4,b_stemp=19.0,b_shum=0,alert=255

type ControlPower bool

var (
	On  = ControlPower(true)
	Off = ControlPower(false)
)

func (j *ControlPower) UnmarshalJSON(b []byte) (err error) {
	*j = ControlPower(bytes.Equal(b, []byte("1")))
	return nil
}

type FanDir int

var (
	Vertical   = FanDir(1)
	Horizontal = FanDir(2)
	Both       = FanDir(3)
)

type FanRate int

func (j *FanRate) UnmarshalJSON(b []byte) (err error) {
	if bytes.Equal(b, []byte(`"A"`)) {
		*j = 1
	} else if bytes.Equal(b, []byte(`"B"`)) {
		*j = 2
	} else {
		var i int
		err = json.Unmarshal(b, &i)
		if err != nil {
			*j = FanRate(i)
		}
	}
	return err
}

func (j *FanRate) encode() string {
	switch *j {
	case 1:
		return "A"
	case 2:
		return "B"
	}
	return strconv.Itoa(int(*j))
}

type ControlInfo struct {
	Power       ControlPower `json:"pow"`
	Mode        int          `json:"mode"`
	SetTemp     float32      `json:"stemp"`
	SetHumidity float32      `json:"shum"`
	FanRate     FanRate      `json:"f_rate"`
	FanDir      FanDir       `json:"f_dir"`
}

type daikinEncode interface {
	forEncode() (v url.Values)
}

func (ci *ControlInfo) forEncode() (v url.Values) {
	v = make(url.Values)

	if ci.Power == On {
		v.Set("pow", "1")
	} else {
		v.Set("pow", "0")
	}
	v.Set("mode", strconv.Itoa(ci.Mode))
	v.Set("stemp", fmt.Sprintf("%.1f", ci.SetTemp))
	v.Set("shum", strconv.Itoa(int(ci.SetHumidity)))
	v.Set("f_rate", ci.FanRate.encode())
	v.Set("f_dir", strconv.Itoa(int(ci.FanDir)))

	return v
}

type BasicInfo struct {
	Version    string        `json:"ver"`
	Name       EncodedString `json:"name"`
	Icon       int           `json:"icon"`
	Method     string        `json:"method"`
	Port       int           `json:"port"`
	GroupName  EncodedString `json:"grp_name"`
	MacAddress string        `json:"mac"`

	// there's a lot more fields here
}

type EncodedString string

func (j *EncodedString) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*j = EncodedString(decodeName(s))
	return nil
}
