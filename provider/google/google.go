/*
 * Copyright 2019 Li Kexian
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

package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/likexian/doh-go/dns"
	"github.com/likexian/gokit/xhttp"
	"github.com/likexian/gokit/xip"
)

// Provider is a DoH provider client
type Provider struct {
	provides int
}

const (
	// DefaultProvides is default provides
	DefaultProvides = iota
)

var (
	// Upstream is DoH query upstream
	Upstream = map[int]string{
		DefaultProvides: "https://dns.google.com/resolve",
	}
)

// Version returns package version
func Version() string {
	return "0.5.3"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// New returns a new google provider client
func New() *Provider {
	return &Provider{
		provides: DefaultProvides,
	}
}

// String returns string of provider
func (c *Provider) String() string {
	return "google"
}

// SetProvides set upstream provides type, google does NOT supported
func (c *Provider) SetProvides(p int) error {
	c.provides = DefaultProvides
	return nil
}

// Query do DoH query
func (c *Provider) Query(ctx context.Context, d dns.Domain, t dns.Type) (*dns.Response, error) {
	return c.ECSQuery(ctx, d, t, "")
}

// ECSQuery do DoH query with the edns0-client-subnet option
func (c *Provider) ECSQuery(ctx context.Context, d dns.Domain, t dns.Type, s dns.ECS) (*dns.Response, error) {
	name, err := d.Punycode()
	if err != nil {
		return nil, err
	}

	param := xhttp.QueryParam{
		"name": name,
		"type": strings.TrimSpace(string(t)),
	}

	ss := strings.TrimSpace(string(s))
	if ss != "" {
		ss, err := xip.FixSubnet(ss)
		if err != nil {
			return nil, err
		}
		param["edns_client_subnet"] = ss
	}

	rsp, err := xhttp.New().Get(ctx, Upstream[c.provides], param, xhttp.Header{"accept": "application/dns-json"})
	if err != nil {
		return nil, err
	}

	defer rsp.Close()
	buf, err := rsp.Bytes()
	if err != nil {
		return nil, err
	}

	rr := &dns.Response{
		Provider: c.String(),
	}
	err = json.NewDecoder(bytes.NewBuffer(buf)).Decode(rr)
	if err != nil {
		return nil, err
	}

	if rr.Status != 0 {
		return rr, fmt.Errorf("doh: google: failed response code %d", rr.Status)
	}

	return rr, nil
}
