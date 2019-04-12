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
 * DNS over HTTPS (DoH) Golang Implementation
 * https://www.likexian.com/
 */

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/likexian/doh-go"
	"github.com/likexian/gokit/xhttp"
	"github.com/likexian/gokit/xip"
	"strings"
)

// Google is a DoH provider
type Google struct {
}

// googleURL is the google DoH url
var googleURL = "https://dns.google.com/resolve"

// String returns string of provider
func (p *Google) String() string {
	return "google"
}

// Query do DoH query
func (p *Google) Query(ctx context.Context, d doh.Domain, t doh.Type) (*doh.Response, error) {
	return p.ECSQuery(ctx, d, t, "")
}

// ECSQuery do DoH query with the edns0-client-subnet option
func (p *Google) ECSQuery(ctx context.Context, d doh.Domain, t doh.Type, s doh.ECS) (*doh.Response, error) {
	param := xhttp.QueryParam{
		"name": strings.TrimSpace(string(d)),
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

	rsp, err := xhttp.Get(googleURL, param, ctx, xhttp.Header{"accept": "application/dns-json"})
	if err != nil {
		return nil, err
	}

	defer rsp.Close()
	buf, err := rsp.Bytes()
	if err != nil {
		return nil, err
	}

	rr := &doh.Response{}
	err = json.NewDecoder(bytes.NewBuffer(buf)).Decode(rr)
	if err != nil {
		return nil, err
	}

	if rr.Status != 0 {
		return rr, fmt.Errorf("doh: failed response code %d", rr.Status)
	}

	return rr, nil
}
