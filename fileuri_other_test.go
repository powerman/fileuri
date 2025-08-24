// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

//go:build !windows

package fileuri_test

var urlTests = []struct {
	url          string
	filePath     string
	canonicalURL string // If empty, assume equal to url.
	wantErr      string
}{
	// --- Absolute paths ---
	// Example from RFC 8089.
	{
		url:      `file:///path/to/file`,
		filePath: `/path/to/file`,
	},
	// Example from RFC 8089.
	{
		url:          `file:/path/to/file`, // Non-canonical single slash
		filePath:     `/path/to/file`,
		canonicalURL: `file:///path/to/file`,
	},

	// --- Root ---
	{
		url:      `file:///`,
		filePath: `/`,
	},

	// --- Special chars ---
	{
		url:      `file:///path/with%20spaces/file.txt`,
		filePath: `/path/with spaces/file.txt`,
	},

	// --- Relative URI (unsupported) ---
	{
		url:     `file:relative.txt`,
		wantErr: `path is not absolute`,
	},

	// --- Remote host (unsupported on POSIX) ---
	{
		url:     `file://server/share/file.txt`,
		wantErr: "file URL specifies non-local host",
	},
	{
		url:     `file://host.example.com/path/to/file`,
		wantErr: "file URL specifies non-local host",
	},

	// --- Extra slashes (RFC forbids, but some parsers allow) ---
	// TODO: Dubious case, maybe should be an error or use path.Clean.
	{
		url:          `file:////path/to/file`,
		filePath:     `//path/to/file`,
		canonicalURL: `file:////path/to/file`,
	},

	// --- Localhost treated as local ---
	// Example from RFC 8089.
	{
		url:          `file://localhost/path/to/file`,
		filePath:     `/path/to/file`,
		canonicalURL: `file:///path/to/file`,
	},

	// --- Missing leading slash ---
	{
		url:     `file://`,
		wantErr: `file URL missing path`,
	},
	{
		url:     `file://localhost`,
		wantErr: `file URL missing path`,
	},

	// --- Wrong scheme ---
	{
		url:     `ftp:///path/to/file`,
		wantErr: `non-file URL`,
	},

	// --- Only file: without slash ---
	{
		url:     `file:`,
		wantErr: `file URL missing path`,
	},
}
