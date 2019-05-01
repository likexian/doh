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

package doh

import (
	"context"
	"github.com/likexian/doh-go/dns"
	"github.com/likexian/gokit/assert"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	assert.Contains(t, Version(), ".")
	assert.Contains(t, Author(), "likexian")
	assert.Contains(t, License(), "Apache License")
}

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := New(CloudflareProvider)
	rsp, err := c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	c = New(DNSPodProvider)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	c = New(GoogleProvider)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)

	c = New(Quad9Provider)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
}

func TestUse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	c := Use()
	defer c.Close()

	_, err := c.Query(ctx, "likexian", dns.TypeA)
	assert.NotNil(t, err)

	c = Use(CloudflareProvider, DNSPodProvider, GoogleProvider, Quad9Provider)
	for i := 0; i < 100; i++ {
		for _, v := range []dns.Type{dns.TypeA, dns.TypeMX} {
			rsp, err := c.Query(ctx, "likexian.com", v)
			assert.Nil(t, err)
			assert.Gt(t, len(rsp.Answer), 0)
		}
	}
}

func TestEnableCache(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	c := Use()
	defer c.Close()

	c.EnableCache(true)
	rsp, err := c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
	ttl := rsp.Answer[0].TTL

	time.Sleep(1 * time.Second)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
	assert.Equal(t, rsp.Answer[0].TTL, ttl)

	c.EnableCache(false)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
	ttl = rsp.Answer[0].TTL

	time.Sleep(1 * time.Second)
	rsp, err = c.Query(ctx, "likexian.com", dns.TypeA)
	assert.Nil(t, err)
	assert.Gt(t, len(rsp.Answer), 0)
	assert.NotEqual(t, rsp.Answer[0].TTL, ttl)
}
