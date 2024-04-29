/*
 * Copyright 2019-2024 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * DNS over HTTPS (DoH) Golang implementation
 * https://www.likexian.com/
 */

package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/likexian/doh/dns"
	"github.com/likexian/gokit/xip"
)

// provider is provider
type provider uint

// Client is DoH provider client
type Client struct {
	provider provider
}

const (
	// DefaultProvider is default provider
	DefaultProvider = iota
	// lastProvider is last provider
	lastProvider
)

var (
	// upstreams is DoH upstreams
	upstreams = map[uint]string{
		DefaultProvider: "https://cloudflare-dns.com/dns-query",
	}
	// httpClient is DoH http client
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 3 * time.Second,
			DisableKeepAlives:   false,
			MaxIdleConns:        256,
			MaxIdleConnsPerHost: 256,
		},
	}
)

// Version returns package version
func Version() string {
	return "0.6.0"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// NewClient returns a new provider client
func NewClient() *Client {
	return &Client{
		provider: DefaultProvider,
	}
}

// String returns string of provider
func (c *Client) String() string {
	return "cloudflare"
}

// SetProvider set upstream provider type, cloudflare does NOT supported
func (c *Client) SetProvider(p provider) error {
	if p >= lastProvider {
		return fmt.Errorf("cloudflare: invalid dns provider")
	}
	c.provider = p
	return nil
}

// Query do DoH query with the edns0-client-subnet option
func (c *Client) Query(ctx context.Context, d dns.Domain, t dns.Type, s ...dns.ECS) (*dns.Response, error) {
	name, err := d.Punycode()
	if err != nil {
		return nil, err
	}

	param := url.Values{}
	param.Add("name", name)
	param.Add("type", strings.TrimSpace(string(t)))

	if len(s) > 0 {
		ss := strings.TrimSpace(string(s[0]))
		if ss != "" {
			ss, err := xip.FixSubnet(ss)
			if err != nil {
				return nil, err
			}
			param.Add("edns_client_subnet", ss)
		}
	}

	dnsURL := fmt.Sprintf("%s?%s", upstreams[uint(c.provider)], param.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dnsURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/dns-json")
	req.Header.Set("User-Agent", fmt.Sprintf("DoH Client/%s", Version()))

	rsp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", rsp.StatusCode)
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	rr := &dns.Response{
		Provider: c.String(),
	}

	err = json.Unmarshal(data, rr)
	if err != nil {
		return nil, err
	}

	if rr.Status != 0 {
		return rr, fmt.Errorf("cloudflare: bad response code: %d", rr.Status)
	}

	return rr, nil
}
