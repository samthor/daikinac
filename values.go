package daikinac

import (
	"bytes"
	"fmt"
	"net/url"
)

func parseValues(b []byte) (v url.Values, err error) {
	v = make(url.Values)

	pairs := bytes.Split(b, []byte(","))
	for _, pair := range pairs {
		parts := bytes.SplitN(pair, []byte("="), 2)
		key := string(parts[0])

		if len(parts) != 2 {
			v.Set(key, "")
			continue
		}

		value := string(parts[1])
		if value == "-" {
			continue
		}
		v.Set(key, value)
	}

	if v.Get("ret") != "OK" {
		return nil, fmt.Errorf("err=%v", v.Get("ret"))
	}

	// log.Printf("got raw output=%+v", v)
	return
}
