package daikinac

import (
	"crypto/tls"
	"net/http"
	"sync"
)

var (
	daikinClient *http.Client

	lockMutex     sync.Mutex
	lockByHostMap map[string]*sync.Mutex
)

func init() {
	// Daikin has a self-signed cert but we don't know what it is so :shrug:
	// we use TLS_RSA_WITH_AES_128_CBC_SHA as it's faster than the default (for my house: 2s -> 1.5s)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		CipherSuites:       []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA},
	}
	daikinClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:    tlsConfig,
			DisableCompression: true,
		},
	}
}

func lockByHost(host string) (m *sync.Mutex) {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	if lockByHostMap == nil {
		lockByHostMap = make(map[string]*sync.Mutex)
	}

	m, ok := lockByHostMap[host]
	if !ok {
		m = &sync.Mutex{}
		lockByHostMap[host] = m
	}
	return m
}
