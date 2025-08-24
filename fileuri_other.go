// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

//go:build !windows

package fileuri

import (
	"errors"
	"path/filepath"
)

var errNonLocalHost = errors.New("file URL specifies non-local host")

func convertFileURLPath(host, path string) (string, error) {
	switch host {
	case "", "localhost":
	default:
		return "", errNonLocalHost
	}
	return filepath.FromSlash(path), nil
}
