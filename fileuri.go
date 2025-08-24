// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package fileuri provides functions to work with file:// URLs.
// It implements RFC 8089 The "file" URI Scheme.
package fileuri

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
)

// TODO(golang.org/issue/32456): If accepted, move these functions into the
// net/url package.

var (
	errNotAbsolute = errors.New("path is not absolute")
	errNonFile     = errors.New("non-file URL")
	errMissingPath = errors.New("file URL missing path")
)

// ToFilePath converts a file URL to a filesystem path.
// It returns an error if the URL is not a file URL or if the path is not absolute.
// The returned path is in the format of the current operating system.
//
// ToFilePath handles file URLs as specified in RFC 8089.
// In particular, it supports both hostless file URLs (file:///path/to/file)
// and UNC file URLs (file://host/share/path/to/file).
//
// ToFilePath does not access the filesystem and does not verify that the path exists.
func ToFilePath(u *url.URL) (string, error) {
	if u.Scheme != "file" {
		return "", errNonFile
	}

	checkAbs := func(path string) (string, error) {
		if !filepath.IsAbs(path) {
			return "", errNotAbsolute
		}
		return path, nil
	}

	if u.Path == "" {
		if u.Host != "" || u.Opaque == "" {
			return "", errMissingPath
		}
		return checkAbs(filepath.FromSlash(u.Opaque))
	}

	path, err := convertFileURLPath(u.Host, u.Path)
	if err != nil {
		return path, err
	}
	return checkAbs(path)
}

// FromFilePath converts an absolute filesystem path to a file URL.
// It returns an error if the path is not absolute.
//
// FromFilePath handles paths on Windows and Unix as specified in RFC 8089.
// In particular, it converts Windows UNC paths (\\host\share\path\to\file)
// to UNC file URLs (file://host/share/path/to/file).
//
// FromFilePath does not access the filesystem and does not verify that the path exists.
func FromFilePath(path string) (*url.URL, error) {
	if !filepath.IsAbs(path) {
		return nil, errNotAbsolute
	}

	// If path has a Windows volume name, convert the volume to a host and prefix
	// per https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/.
	if vol := filepath.VolumeName(path); vol != "" {
		if strings.HasPrefix(vol, `\\`) {
			path = filepath.ToSlash(path[2:])
			i := strings.IndexByte(path, '/')

			if i < 0 {
				// A degenerate case.
				// \\host.example.com (without a share name)
				// becomes
				// file://host.example.com/
				return &url.URL{
					Scheme: "file",
					Host:   path,
					Path:   "/",
				}, nil
			}

			// \\host.example.com\Share\path\to\file
			// becomes
			// file://host.example.com/Share/path/to/file
			return &url.URL{
				Scheme: "file",
				Host:   path[:i],
				Path:   filepath.ToSlash(path[i:]),
			}, nil
		}

		// C:\path\to\file
		// becomes
		// file:///C:/path/to/file
		return &url.URL{
			Scheme: "file",
			Path:   "/" + filepath.ToSlash(path),
		}, nil
	}

	// /path/to/file
	// becomes
	// file:///path/to/file
	return &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(path),
	}, nil
}
