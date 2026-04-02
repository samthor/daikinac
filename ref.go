package daikinac

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Device is a split-system provided by Daikin.
//
// If UUID is set, this is a new-style device which requires key/UUID registration.
// It uses (pointless) HTTPS, which slows down access to the device.
// However, doing calls in parallel is faster than solo, so consider using a WaitGroup.
type Device struct {
	Host string // IP or hostname
	UUID string // UUID from modern devices (blank otherwise): part of registration flow
}

// Do performs a request to the device.
// If out is non-nil, the output keys/values will be encoded there "as JSON".
//
// You may want to use DoSetControl, DoGetControl, or DoGetSensor for type-safe calls.
func (d *Device) Do(c context.Context, p string, in, out any) (err error) {
	protocol := "http"

	if d.UUID != "" {
		protocol = "https"
	} else {
		// hilariously, the old devices struggle with more than one req - so we serialize them per-host.
		// this is super lazy and prevents Device from containing the lock
		lock := lockByHost(d.Host)
		lock.Lock()
		defer lock.Unlock()
	}

	u, err := url.Parse(fmt.Sprintf("%s://%s%s", protocol, d.Host, p))
	if err != nil {
		return err
	}
	method := http.MethodGet

	// pass values
	if in != nil {
		fe, ok := in.(daikinEncode)
		if !ok {
			return fmt.Errorf("cannot send unknown type")
		}
		values := fe.asValues()
		u.RawQuery = values.Encode()
	}

	// do thing!
	request, err := http.NewRequestWithContext(c, method, u.String(), http.NoBody)
	if err != nil {
		return err
	}
	if d.UUID != "" {
		request.Header["X-Daikin-uuid"] = []string{d.UUID} // set directly as Daikin doesn't respect normalization
	}

	resp, err := daikinClient.Do(request) // this takes all the time (obvs)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	v, err := parseValues(b)
	if err != nil || out == nil {
		return err
	}

	fd, ok := out.(daikinDecode)
	if !ok {
		return fmt.Errorf("can't parse into unknown type")
	}
	return fd.fromValues(v)
}

func (d *Device) KeepAlive(ctx context.Context) (err error) {
	if d.UUID == "" {
		return
	}

	u, err := url.Parse(fmt.Sprintf("https://%s/common/get_basic_info", d.Host))
	if err != nil {
		return err
	}

	// do thing!
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return err
	}
	if d.UUID != "" {
		request.Header["X-Daikin-uuid"] = []string{d.UUID} // set directly as Daikin doesn't respect normalization
	}

	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		default:
		}

		log.Printf("KEEPALIVE")
		resp, err := daikinClient.Do(request) // this takes all the time (obvs)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
}
