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

package doh

import (
	"context"
	"github.com/likexian/doh-go/dns"
	"github.com/likexian/doh-go/provider/cloudflare"
	"github.com/likexian/doh-go/provider/dnspod"
	"github.com/likexian/doh-go/provider/google"
	"github.com/likexian/doh-go/provider/quad9"
)

// Provider is the provider interface
type Provider interface {
	Query(context.Context, dns.Domain, dns.Type) (*dns.Response, error)
	ECSQuery(context.Context, dns.Domain, dns.Type, dns.ECS) (*dns.Response, error)
	String() string
}

// DoH Providers
const (
	CloudflareProvider = iota
	DNSPodProvider
	GoogleProvider
	Quad9Provider
)

// Version returns package version
func Version() string {
	return "0.3.0"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// New returns a new DoH client, quad9 is default
func New(provider int) Provider {
	switch provider {
	case CloudflareProvider:
		return cloudflare.New()
	case DNSPodProvider:
		return dnspod.New()
	case GoogleProvider:
		return google.New()
	default:
		return quad9.New()
	}
}
