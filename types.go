package daikinac

import (
	"fmt"
	"net/url"
	"strconv"
)

type daikinEncode interface {
	asValues() (v url.Values)
}

type daikinDecode interface {
	fromValues(v url.Values) (err error)
}

type SensorInfo struct {
	HomeTemp       float64
	HomeHumidity   int
	OutsideTemp    float64
	CompressorFreq int
}

func (si *SensorInfo) fromValues(v url.Values) (err error) {
	si.HomeTemp, _ = strconv.ParseFloat(v.Get("htemp"), 64)
	si.HomeHumidity, _ = strconv.Atoi(v.Get("hhum"))
	si.OutsideTemp, _ = strconv.ParseFloat(v.Get("otemp"), 64)
	si.CompressorFreq, _ = strconv.Atoi(v.Get("cmp"))
	return nil
}

// ret=OK,pow=0,mode=4,adv=,stemp=19.0,shum=0,dt1=25.0,dt2=M,dt3=25.0,dt4=19.0,dt5=19.0,dt7=25.0,dh1=AUTO,dh2=50,dh3=0,dh4=0,dh5=0,dh7=AUTO,dhh=50,b_mode=4,b_stemp=19.0,b_shum=0,alert=255

type ControlPower bool

var (
	On  = ControlPower(true)
	Off = ControlPower(false)
)

type FanDir int

var (
	FanNone       = FanDir(0)
	FanVertical   = FanDir(1)
	FanHorizontal = FanDir(2)
	FanBoth       = FanDir(3)
)

type FanRate int

var (
	FanAuto  = FanRate(1)
	FanQuiet = FanRate(2)
)

func (j *FanRate) decode(x string) (err error) {
	switch x {
	case "A":
		*j = 1
	case "B":
		*j = 2
	default:
		var v int
		v, err = strconv.Atoi(x)
		*j = FanRate(v)
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

type Mode int

var (
	ModeAuto = Mode(0)
	ModeDry  = Mode(2)
	ModeCool = Mode(3)
	ModeHeat = Mode(4)
	ModeFan  = Mode(6)
)

type ControlInfo struct {
	Power ControlPower
	PrimaryControl

	PriorModes []ControlInfoMode // ignored for encoding, "H" mode is put into 0
	BMode      PrimaryControl    // ignored for encoding
}

type PrimaryControl struct {
	Mode Mode
	ControlInfoMode
}

func (pc *PrimaryControl) fromValues(v url.Values, prefix string) {
	get := func(key string) string { return v.Get(fmt.Sprintf("%s%s", prefix, key)) }

	mode, _ := strconv.Atoi(get("mode"))
	pc.Mode = Mode(mode)

	pc.SetTemp, _ = strconv.ParseFloat(get("stemp"), 64)
	pc.SetHumidity = parseHumidity(get("shum"))
	pc.FanRate.decode(get("f_rate"))

	fanDir, _ := strconv.Atoi(get("f_dir"))
	pc.FanDir = FanDir(fanDir)
}

type ControlInfoMode struct {
	SetTemp     float64
	SetHumidity int
	FanRate     FanRate
	FanDir      FanDir
}

func (ci *ControlInfo) asValues() (v url.Values) {
	v = make(url.Values)

	if ci.Power == On {
		v.Set("pow", "1")
	} else {
		v.Set("pow", "0")
	}
	v.Set("mode", strconv.Itoa(int(ci.Mode)))
	v.Set("stemp", fmt.Sprintf("%.1f", ci.SetTemp))
	v.Set("shum", renderHumidity(ci.SetHumidity))
	v.Set("f_rate", ci.FanRate.encode())
	v.Set("f_dir", strconv.Itoa(int(ci.FanDir)))

	return v
}

func parseHumidity(s string) (out int) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return v
}

func renderHumidity(v int) (out string) {
	if v >= 0 && v <= 100 {
		return strconv.Itoa(v)
	}
	return "AUTO"
}

func (ci *ControlInfo) fromValues(v url.Values) (err error) {
	ci.Power = v.Get("pow") == "1"

	ci.PrimaryControl.fromValues(v, "")
	ci.BMode.fromValues(v, "b_")

	// parse prior modes
	for x := range 8 {
		suffix := "h"
		if x > 0 {
			suffix = strconv.Itoa(x)
		}
		get := func(key string) string { return v.Get(fmt.Sprintf("%s%s", key, suffix)) }

		var prior ControlInfoMode

		prior.SetTemp, _ = strconv.ParseFloat(get("dt"), 64)
		prior.SetHumidity = parseHumidity(get("dh"))
		prior.FanRate.decode(get("dfr"))

		fanDir, _ := strconv.Atoi(get("dfd"))
		prior.FanDir = FanDir(fanDir)

		ci.PriorModes = append(ci.PriorModes, prior)
	}

	return err
}

type BasicInfo struct {
	Version    string
	Name       string
	Icon       int
	Method     string
	Port       int
	GroupName  string
	MacAddress string

	// there's a lot more fields here
}

func (bi *BasicInfo) fromValues(v url.Values) (err error) {
	bi.Version = v.Get("ver")
	bi.Name = decodeName(v.Get("name"))
	bi.Icon, _ = strconv.Atoi(v.Get("icon"))
	bi.Method = v.Get("method")
	bi.Port, _ = strconv.Atoi(v.Get("port"))
	bi.GroupName = decodeName(v.Get("grp_name"))
	bi.MacAddress = v.Get("mac")
	return nil
}
