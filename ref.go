package daikinac

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Device is a split-system provided by Daikin.
type Device struct {
	Host string // IP or hostname
	UUID string // UUID from modern devices: part of registration flow
}

// Do performs a request to the device.
// If out is non-nil, the output keys/values will be encoded there "as JSON".
func (d *Device) Do(c context.Context, p string, in, out any) (err error) {
	protocol := "http"
	client := http.DefaultClient

	if d.UUID != "" {
		// Daikin has a self-signed cert but we don't know what it is so :shrug:
		protocol = "https"
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true, Renegotiation: tls.RenegotiateFreelyAsClient},
			},
		}
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
	request.Header["X-Daikin-uuid"] = []string{d.UUID} // set directly as Daikin doesn't respect normalization

	resp, err := client.Do(request)
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
