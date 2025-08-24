# fileuri is a Go package to work with file:// URI

[![License MIT](https://img.shields.io/badge/license-MIT-royalblue.svg)](LICENSE)
[![Go version](https://img.shields.io/github/go-mod/go-version/powerman/fileuri?color=blue)](https://go.dev/)
[![Test](https://img.shields.io/github/actions/workflow/status/powerman/fileuri/test.yml?label=test)](https://github.com/powerman/fileuri/actions/workflows/test.yml)
[![Coverage Status](https://raw.githubusercontent.com/powerman/fileuri/gh-badges/coverage.svg)](https://github.com/powerman/fileuri/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/fileuri)](https://goreportcard.com/report/github.com/powerman/fileuri)
[![Release](https://img.shields.io/github/v/release/powerman/fileuri?color=blue)](https://github.com/powerman/fileuri/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/fileuri.svg)](https://pkg.go.dev/github.com/powerman/fileuri)

![Linux | amd64 arm64 armv7 ppc64le s390x riscv64](https://img.shields.io/badge/Linux-amd64%20arm64%20armv7%20ppc64le%20s390x%20riscv64-royalblue)
![macOS | amd64 arm64](https://img.shields.io/badge/macOS-amd64%20arm64-royalblue)
![Windows | amd64 arm64](https://img.shields.io/badge/Windows-amd64%20arm64-royalblue)

Implements [RFC 8089 The "file" URI Scheme](https://datatracker.ietf.org/doc/html/rfc8089).

The implementation is based on `cmd/go/internal/web` package from Go 1.25.0.

See also [Go issue 32456](https://github.com/golang/go/issues/32456).
