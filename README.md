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

## Features

- DoH client, Simple and Easy to use
- Support cloudflare, google, quad9 and dnspod
- Specify the provider you like
- Auto select fastest provider
- Enable cache is supported
- EDNS0-Client-Subnet query supported

## Installation

    go get -u github.com/likexian/doh-go

## Importing

    import (
        "github.com/likexian/doh-go"
        "github.com/likexian/doh-go/dns"
    )

## Documentation

Visit the docs on [GoDoc](https://godoc.org/github.com/likexian/doh-go)

## Example

### Select fastest provider and query (Highly Recommend)

```go
// init a context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// init doh client, auto select the fastest provider base on your like
// you can also use as: c := doh.Use(), it will select from all providers
c := doh.Use(doh.CloudflareProvider, doh.GoogleProvider)

// do doh query
rsp, err := c.Query(ctx, "likexian.com", dns.TypeA)
if err != nil {
    panic(err)
}

// close the client
c.Close()

// doh dns answer
answer := rsp.Answer

// print all answer
for _, a := range answer {
    fmt.Printf("%s -> %s\n", a.Name, a.Data)
}
```

### Specify DoH provider and query (You are Welcome)

```go
// init a context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// init doh client, specify one provider
c := doh.New(Quad9Provider)

// do doh query
rsp, err := c.Query(ctx, "likexian.com", dns.TypeMX)
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

### Google (NOT work in Mainland China)

Google Public DNS is a recursive DNS resolver, similar to other publicly available services. We think it provides many benefits, including improved security, fast performance, and more valid results. But it is not work in mainland China.

- https://developers.google.com/speed/public-dns/docs/dns-over-https

### DNSPod (Fake DoH)

DNS over HTTP but NOT HTTPS and A record only. This is something known as HTTPDNS, provided by DNSPod (Tencent Cloud). The backend is a anycast public DNS platform well known in China.

- https://cloud.tencent.com/document/product/379/3524

## LICENSE

Copyright 2019 Li Kexian

Licensed under the Apache License 2.0

## About

- [Li Kexian](https://www.likexian.com/)

## DONATE

- [Help me make perfect](https://www.likexian.com/donate/)
