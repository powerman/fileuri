// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package fileuri

import (
	"errors"
	"path/filepath"
	"strings"
)

const maxPath = 260 // Windows MAX_PATH limitation.

var (
	errVolumeInHost  = errors.New("file URL encodes volume in host field: too few slashes?")
	errMissingVolume = errors.New("file URL missing drive letter")
	errMissingShare  = errors.New("file URL missing UNC share name")
)

func convertFileURLPath(host, path string) (string, error) {
	if path == "" || path[0] != '/' {
		return "", errNotAbsolute
	}

	path = filepath.FromSlash(path)

	// We interpret Windows file URLs per the description in
	// https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/.

	// The host part of a file URL (if any) is the UNC volume name,
	// but RFC 8089 reserves the authority "localhost" for the local machine.
	if host != "" && host != "localhost" {
		// A common "legacy" format omits the leading slash before a drive letter,
		// encoding the drive letter as the host instead of part of the path.
		// (See https://blogs.msdn.microsoft.com/freeassociations/2005/05/19/the-bizarre-and-unhappy-story-of-file-urls/.)
		// We do not support that format, but we should at least emit a more
		// helpful error message for it.
		if filepath.VolumeName(host) != "" {
			return "", errVolumeInHost
		}
		if path == `\` {
			return "", errMissingShare
		}
		path = `\\` + host + path
		if len(path) >= maxPath { // Too strict because check bytes instead of UTF-16 chars.
			// Use the extended-length path prefix to avoid MAX_PATH limitations.
			// See https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#maximum-path-length-limitation.
			path = `\\?\UNC` + filepath.Clean(path[1:])
		}
		return path, nil
	}

	// If host is empty, path must contain an initial slash followed by a
	// drive letter and path. Remove the slash and verify that the path is valid.
	if vol := filepath.VolumeName(path[1:]); vol == "" || strings.HasPrefix(vol, `\\`) {
		return "", errMissingVolume
	}
	path = path[1:]
	if len(path) >= maxPath { // Too strict because check bytes instead of UTF-16 chars.
		path = `\\?\` + filepath.Clean(path)
	}
	return path, nil
}
