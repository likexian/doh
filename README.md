# doh.go

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/likexian/doh-go?status.svg)](https://godoc.org/github.com/likexian/doh-go)
[![Build Status](https://travis-ci.org/likexian/doh-go.svg?branch=master)](https://travis-ci.org/likexian/doh-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/likexian/doh-go)](https://goreportcard.com/report/github.com/likexian/doh-go)
[![Code Cover](https://codecov.io/gh/likexian/doh-go/graph/badge.svg)](https://codecov.io/gh/likexian/doh-go)

doh-go is a DNS over HTTPS (DoH) Golang client implementation.

## Overview

DNS over HTTPS (DoH) is a protocol for performing remote Domain Name System (DNS) resolution via the HTTPS protocol. Specification is [RFC 8484 - DNS Queries over HTTPS (DoH)](https://tools.ietf.org/html/rfc8484).

This module provides a easy way to using DoH as client in golang.

## Installation

    go get -u github.com/likexian/doh-go

## Importing

    import (
        "github.com/likexian/doh-go"
    )

## Documentation

Visit the docs on [GoDoc](https://godoc.org/github.com/likexian/doh-go)

## Example

```go
// init a context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// init doh client
c := doh.New(Quad9Provider)

// do doh query
rsp, err := c.Query(ctx, "likexian.com", dns.TypeA)
if err != nil {
    panic(err)
}

// doh dns answer
answer := rsp.Answer

// print all answer
for _, a := range answer {
    fmt.Printf("%s -> %s\n", a.Name, a.Data)
}
```

## Providers

### Quad9 (Recommend)

Quad9 is a free, recursive, anycast DNS platform that provides end users robust security protections, high-performance, and privacy.

- https://www.quad9.net/doh-quad9-dns-servers/

### Cloudflare (Fast)

Cloudflare's mission is to help build a better Internet. We're excited today to take another step toward that mission with the launch of 1.1.1.1 — the Internet's fastest, privacy-first consumer DNS service.

- https://developers.cloudflare.com/1.1.1.1/dns-over-https/

### Google (NOT work in China mainlan)

Google Public DNS is a recursive DNS resolver, similar to other publicly available services. We think it provides many benefits, including improved security, fast performance, and more valid results.

- https://developers.google.com/speed/public-dns/docs/dns-over-https

## LICENSE

Copyright 2019 Li Kexian

Licensed under the Apache License 2.0

## About

- [Li Kexian](https://www.likexian.com/)

## DONATE

- [Help me make perfect](https://www.likexian.com/donate/)
