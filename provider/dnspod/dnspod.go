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

package dnspod

import (
	"context"
	"fmt"
	"strconv"
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
		DefaultProvides: "http://119.29.29.29/d",
	}
)

// Version returns package version
func Version() string {
	return "0.2.4"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// New returns a new dnspod provider client
func New() *Provider {
	return &Provider{
		provides: DefaultProvides,
	}
}

// String returns string of provider
func (c *Provider) String() string {
	return "dnspod"
}

// SetProvides set upstream provides type, dnspod does NOT supported
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
	if t != dns.TypeA {
		return nil, fmt.Errorf("doh: dnspod: only A record type is supported")
	}

	name, err := d.Punycode()
	if err != nil {
		return nil, err
	}

	param := xhttp.QueryParam{
		"dn":  name,
		"ttl": "1",
	}

	ss := strings.TrimSpace(string(s))
	if ss != "" {
		ss, err := xip.FixSubnet(ss)
		if err != nil {
			return nil, err
		}
		ips := strings.Split(ss, "/")
		param["ip"] = ips[0]
	}

	rsp, err := xhttp.New().Get(ctx, Upstream[c.provides], param)
	if err != nil {
		return nil, err
	}

	defer rsp.Close()
	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("doh: dnspod: bad status code: %d", rsp.StatusCode)
	}

	txt, err := rsp.String()
	if err != nil {
		return nil, err
	}

	rr := &dns.Response{
		Status:   0,
		TC:       false,
		RD:       true,
		RA:       true,
		AD:       false,
		CD:       false,
		Question: []dns.Question{},
		Answer:   []dns.Answer{},
		Provider: c.String(),
	}
	rr.Question = append(rr.Question, dns.Question{Name: name, Type: 1})

	txt = strings.TrimSpace(txt)
	if txt == "" {
		rr.Status = 3
		return rr, fmt.Errorf("doh: dnspod: empty response from server")
	}

	ttl := 0
	ts := strings.Split(txt, ",")
	if len(ts) == 2 {
		i, err := strconv.Atoi(ts[1])
		if err == nil {
			ttl = i
		}
	}

	ts = strings.Split(ts[0], ";")
	for _, v := range ts {
		if xip.IsIP(v) {
			rr.Answer = append(rr.Answer, dns.Answer{Name: name, Type: 1, TTL: ttl, Data: v})
		}
	}

	return rr, nil
}
