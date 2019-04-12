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
	"github.com/likexian/gokit/xhttp"
)

// Client is a DoH provider client
type Client struct {
	xhttp *xhttp.Request
}

// Provider is the provider interface
type Provider interface {
	New(...string) *Client
	Query(context.Context, Domain, Type) (*Response, error)
	ECSQuery(context.Context, Domain, Type, ECS) (*Response, error)
	String() string
}
