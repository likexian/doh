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

package doh

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/likexian/doh/dns"
	"github.com/likexian/doh/provider/cloudflare"
	"github.com/likexian/doh/provider/dnspod"
	"github.com/likexian/doh/provider/google"
	"github.com/likexian/doh/provider/quad9"
	"github.com/likexian/gokit/xcache"
	"github.com/likexian/gokit/xhash"
)

// provider is provider
type provider uint

// Provider is the provider interface
type Provider interface {
	Query(context.Context, dns.Domain, dns.Type, ...dns.ECS) (*dns.Response, error)
	String() string
}

// DoH is doh client
type DoH struct {
	providers []Provider
	cache     xcache.Cachex
	stats     map[int][]interface{}
	stopc     chan bool
	sync.RWMutex
}

// DoH Providers enum
const (
	CloudflareProvider provider = iota
	DNSPodProvider
	GoogleProvider
	Quad9Provider
)

// DoH Providers list
var (
	Providers = []provider{
		CloudflareProvider,
		DNSPodProvider,
		GoogleProvider,
		Quad9Provider,
	}
)

// Version returns package version
func Version() string {
	return "0.7.1"
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
func New(provider provider) Provider {
	switch provider {
	case CloudflareProvider:
		return cloudflare.NewClient()
	case DNSPodProvider:
		return dnspod.NewClient()
	case GoogleProvider:
		return google.NewClient()
	default:
		return quad9.NewClient()
	}
}

// Use returns a new DoH client,
// You can specify one or multiple provider,
// if multiple, it will try to select the fastest
func Use(provider ...provider) *DoH {
	c := &DoH{
		providers: []Provider{},
		cache:     nil,
		stats:     map[int][]interface{}{},
		stopc:     make(chan bool),
	}

	if len(provider) == 0 {
		provider = Providers
	}

	for _, v := range provider {
		c.providers = append(c.providers, New(v))
	}

	go func() {
		t := time.NewTicker(time.Duration(5) * time.Second)
		for {
			select {
			case <-c.stopc:
				t.Stop()
				return
			case <-t.C:
				c.Lock()
				c.stats = map[int][]interface{}{}
				c.Unlock()
			}
		}
	}()

	return c
}

// EnableCache enable query cache
func (c *DoH) EnableCache(cache bool) *DoH {
	if cache {
		c.cache = xcache.New(xcache.MemoryCache)
	} else {
		c.cache = nil
	}

	return c
}

// Close close doh client
func (c *DoH) Close() {
	c.stopc <- true
	if c.cache != nil {
		c.cache.Close()
	}
}

// Query do DoH query
func (c *DoH) Query(ctx context.Context, d dns.Domain, t dns.Type, s ...dns.ECS) (*dns.Response, error) {
	providers := c.providers

	c.RLock()
	if len(c.stats) > 0 {
		min := []interface{}{0, 100.0}
		for k, v := range c.stats {
			r := v[2].(float64)
			if r < min[1].(float64) {
				min = []interface{}{k, r}
			}
		}
		providers = []Provider{c.providers[min[0].(int)]}
	}
	c.RUnlock()

	return c.fastQuery(ctx, providers, d, t, s...)
}

// fastQuery do query and returns the fastest result
func (c *DoH) fastQuery(ctx context.Context,
	ps []Provider, d dns.Domain, t dns.Type, s ...dns.ECS) (*dns.Response, error) {
	cacheKey := ""
	if c.cache != nil {
		var ss string
		if len(s) > 0 && s[0] != "" {
			ss = strings.TrimSpace(string(s[0]))
		}
		cacheKey = xhash.Sha1(string(d), string(t), ss).Hex()
		v := c.cache.Get(cacheKey)
		if v != nil {
			return v.(*dns.Response), nil
		}
	}

	ctxs, cancels := context.WithCancel(ctx)
	defer cancels()

	r := make(chan interface{})
	for k, p := range ps {
		go func(k int, p Provider) {
			rsp, err := p.Query(ctxs, d, t, s...)
			c.Lock()
			if _, ok := c.stats[k]; !ok {
				c.stats[k] = []interface{}{0, 0, 100}
			}
			c.stats[k][1] = c.stats[k][1].(int) + 1
			if err != nil {
				c.stats[k][0] = c.stats[k][0].(int) + 1
			}
			c.stats[k][2] = float64(c.stats[k][0].(int)) / float64(c.stats[k][1].(int))
			c.Unlock()
			if err == nil {
				r <- rsp
			} else {
				r <- nil
			}
		}(k, p)
	}

	total := 0
	result := &dns.Response{
		Status: -1,
	}

	for v := range r {
		total++
		if v != nil {
			cancels()
			result = v.(*dns.Response)
			if cacheKey != "" {
				ttl := 30
				if len(result.Answer) > 0 {
					ttl = result.Answer[0].TTL
				}
				_ = c.cache.Set(cacheKey, result, int64(ttl))
			}
		}
		if total >= len(ps) {
			close(r)
			break
		}
	}

	if result.Status == -1 {
		return nil, fmt.Errorf("doh: all query failed")
	}

	return result, nil
}
