// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package fileuri_test

import (
	"net/url"
	"testing"

	"github.com/powerman/fileuri"
)

func TestURLToFilePath(t *testing.T) {
	for _, tc := range urlTests {
		if tc.url == "" {
			continue
		}

		t.Run(tc.url, func(t *testing.T) {
			u, err := url.Parse(tc.url)
			if err != nil {
				t.Fatalf("url.Parse(%q): %v", tc.url, err)
			}

			path, err := fileuri.ToFilePath(u)
			if err != nil {
				if err.Error() == tc.wantErr {
					return
				}
				if tc.wantErr == "" {
					t.Fatalf("ToFilePath(%v): %v; want <nil>", u, err)
				} else {
					t.Fatalf("ToFilePath(%v): %v; want %s", u, err, tc.wantErr)
				}
			}

			if path != tc.filePath || tc.wantErr != "" {
				t.Fatalf("ToFilePath(%v) = %q, <nil>; want %q, %s", u, path, tc.filePath, tc.wantErr)
			}
		})
	}
}

func TestURLFromFilePath(t *testing.T) {
	for _, tc := range urlTests {
		if tc.filePath == "" {
			continue
		}

		t.Run(tc.filePath, func(t *testing.T) {
			u, err := fileuri.FromFilePath(tc.filePath)
			if err != nil {
				if err.Error() == tc.wantErr {
					return
				}
				if tc.wantErr == "" {
					t.Fatalf("FromFilePath(%v): %v; want <nil>", tc.filePath, err)
				} else {
					t.Fatalf("FromFilePath(%v): %v; want %s", tc.filePath, err, tc.wantErr)
				}
			}

			if tc.wantErr != "" {
				t.Fatalf("FromFilePath(%v) = <nil>; want error: %s", tc.filePath, tc.wantErr)
			}

			wantURL := tc.url
			if tc.canonicalURL != "" {
				wantURL = tc.canonicalURL
			}
			if u.String() != wantURL {
				t.Errorf("FromFilePath(%v) = %v; want %s", tc.filePath, u, wantURL)
			}
		})
	}
}
