package daikinac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func parseValues(b []byte, target any) (err error) {
	tmp := make(map[string]any)

	pairs := bytes.Split(b, []byte(","))
	for _, pair := range pairs {
		parts := bytes.SplitN(pair, []byte("="), 2)
		key := string(parts[0])

		if len(parts) != 2 {
			tmp[key] = ""
			continue
		}

		value := string(parts[1])
		if value == "-" {
			continue
		} else if strings.ContainsRune(value, '_') {
			// ok
		} else if f, err := strconv.ParseFloat(value, 32); err == nil {
			tmp[key] = f
			continue
		} else {
			// ok
		}
		tmp[key] = value
	}

	if tmp["ret"] != "OK" {
		return fmt.Errorf("err=%v", tmp["ret"])
	}

	btmp, _ := json.Marshal(tmp)
	if target != nil {
		return json.Unmarshal(btmp, target)
	}
	return nil
}
