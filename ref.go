package daikinac

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

/*
  const daikinConfig: Record<string, { name: string; at: string }> = {
    bedroom: { name: 'Bedroom', at: '192.168.3.204' },
    loft: { name: 'Loft', at: '192.168.3.225' },
    'living-room': { name: 'Living Room', at: '192.168.3.152' },
    den: { name: 'Den', at: '192.168.3.146' },
    office: { name: 'Office', at: '192.168.3.245/f45aab28604811eca7c4737954d1686f' },
  };
*/

type Device struct {
	Host string
	UUID string
}

func (d *Device) Do(c context.Context, p string, in, out any) (err error) {
	protocol := "http"
	client := http.DefaultClient

	if d.UUID != "" {
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

	if in != nil {
		fe, ok := in.(daikinEncode)
		if !ok {
			return fmt.Errorf("cannot send unknown type")
		}
		values := fe.forEncode()

		// TODO: only works with UUID?
		u.RawQuery = values.Encode()
	}

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

	return parseValues(b, out)
}
