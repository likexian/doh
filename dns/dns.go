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

package dns

import (
	"strings"

	"golang.org/x/net/idna"
)

// Domain is dns query domain
type Domain string

// Type is dns query type
type Type string

// ECS is the edns0-client-subnet option, for example: 1.2.3.4/24
type ECS string

// Question is dns query question
type Question struct {
	Name string `json:"name"`
	Type int    `json:"type"`
}

// Answer is dns query answer
type Answer struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// Response is dns query response
type Response struct {
	Status   int        `json:"Status"`
	TC       bool       `json:"TC"`
	RD       bool       `json:"RD"`
	RA       bool       `json:"RA"`
	AD       bool       `json:"AD"`
	CD       bool       `json:"CD"`
	Question []Question `json:"Question"`
	Answer   []Answer   `json:"Answer"`
	Provider string     `json:"provider"`
}

// Supported dns query type
var (
	TypeA     = Type("A")
	TypeAAAA  = Type("AAAA")
	TypeCNAME = Type("CNAME")
	TypeMX    = Type("MX")
	TypeTXT   = Type("TXT")
	TypeSPF   = Type("SPF")
	TypeNS    = Type("NS")
	TypeSOA   = Type("SOA")
	TypePTR   = Type("PTR")
	TypeANY   = Type("ANY")
)

// Version returns package version
func Version() string {
	return "0.3.2"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// Punycode returns punycode of domain
func (d Domain) Punycode() (string, error) {
	name := strings.TrimSpace(string(d))

	return idna.New(
		idna.MapForLookup(),
		idna.Transitional(true),
		idna.StrictDomainName(false),
	).ToASCII(name)
}
