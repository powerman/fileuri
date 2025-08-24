// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Copyright 2025 Alex Efros <powerman@powerman.name>. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

package fileuri_test

import "strings"

var urlTests = []struct {
	url          string
	filePath     string
	canonicalURL string // If empty, assume equal to url.
	wantErr      string
}{
	// Examples from https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/:

	{
		url:      `file://laptop/My%20Documents/FileSchemeURIs.doc`,
		filePath: `\\laptop\My Documents\FileSchemeURIs.doc`,
	},
	{
		url:      `file:///C:/Documents%20and%20Settings/davris/FileSchemeURIs.doc`,
		filePath: `C:\Documents and Settings\davris\FileSchemeURIs.doc`,
	},
	{
		url:      `file:///D:/Program%20Files/Viewer/startup.htm`,
		filePath: `D:\Program Files\Viewer\startup.htm`,
	},
	{
		url:          `file:///C:/Program%20Files/Music/Web%20Sys/main.html?REQUEST=RADIO`,
		filePath:     `C:\Program Files\Music\Web Sys\main.html`,
		canonicalURL: `file:///C:/Program%20Files/Music/Web%20Sys/main.html`,
	},
	{
		url:      `file://applib/products/a-b/abc_9/4148.920a/media/start.swf`,
		filePath: `\\applib\products\a-b\abc_9\4148.920a\media\start.swf`,
	},
	{
		url:     `file:////applib/products/a%2Db/abc%5F9/4148.920a/media/start.swf`,
		wantErr: "file URL missing drive letter",
	},
	{
		url:     `C:\Program Files\Music\Web Sys\main.html?REQUEST=RADIO`,
		wantErr: "non-file URL",
	},

	// The example "file://D:\Program Files\Viewer\startup.htm" errors out in
	// url.Parse, so we substitute a slash-based path for testing instead.
	{
		url:     `file://D:/Program Files/Viewer/startup.htm`,
		wantErr: "file URL encodes volume in host field: too few slashes?",
	},

	// The blog post discourages the use of non-ASCII characters because they
	// depend on the user's current codepage. However, when we are working with Go
	// strings we assume UTF-8 encoding, and our url package refuses to encode
	// URLs to non-ASCII strings.
	{
		url:          `file:///C:/exampleㄓ.txt`,
		filePath:     `C:\exampleㄓ.txt`,
		canonicalURL: `file:///C:/example%E3%84%93.txt`,
	},
	{
		url:      `file:///C:/example%E3%84%93.txt`,
		filePath: `C:\exampleㄓ.txt`,
	},

	// Examples from RFC 8089:

	// We allow the drive-letter variation from section E.2, because it is
	// simpler to support than not to. However, we do not generate the shorter
	// form in the reverse direction.
	{
		url:          `file:c:/path/to/file`,
		filePath:     `c:\path\to\file`,
		canonicalURL: `file:///c:/path/to/file`,
	},

	// We encode the UNC share name as the authority following section E.3.1,
	// because that is what the Microsoft blog post explicitly recommends.
	{
		url:      `file://host.example.com/Share/path/to/file.txt`,
		filePath: `\\host.example.com\Share\path\to\file.txt`,
	},

	// We decline the four- and five-slash variations from section E.3.2.
	// The paths in these URLs would change meaning under path.Clean.
	{
		url:     `file:////host.example.com/path/to/file`,
		wantErr: "file URL missing drive letter",
	},
	{
		url:     `file://///host.example.com/path/to/file`,
		wantErr: "file URL missing drive letter",
	},

	// --- Absolute disk paths ---
	{
		url:      `file:///C:/path/to/file.txt`,
		filePath: `C:\path\to\file.txt`,
	},
	{
		url:          `file:/C:/path/to/file.txt`, // Non-canonical
		filePath:     `C:\path\to\file.txt`,
		canonicalURL: `file:///C:/path/to/file.txt`,
	},

	// --- Root of drive ---
	{
		url:      `file:///C:/`,
		filePath: `C:\`,
	},

	// --- Special chars ---
	{
		url:      `file:///C:/Program%20Files/App/app.exe`,
		filePath: `C:\Program Files\App\app.exe`,
	},

	// --- Relative URI (unsupported) ---
	{
		url:     `file:relative.txt`,
		wantErr: `path is not absolute`,
	},

	// --- UNC standard ---
	{
		url:      `file://server/share/file.txt`,
		filePath: `\\server\share\file.txt`,
	},

	// --- UNC with extra slashes ---
	{
		url:     `file:////server/share/file.txt`,
		wantErr: `file URL missing drive letter`,
	},
	{
		url:     `file://///server/share/file.txt`,
		wantErr: `file URL missing drive letter`,
	},

	// --- Localhost treated as disk path ---
	{
		url:          `file://localhost/C:/Windows/System32/cmd.exe`,
		filePath:     `C:\Windows\System32\cmd.exe`,
		canonicalURL: `file:///C:/Windows/System32/cmd.exe`,
	},

	// --- Extended-length disk path (normalize) ---
	// Path longer than MAX_PATH (260 characters including NUL) requires
	// extended-length prefix \\?\.
	// TODO: Characters are UTF-16, U+10000 … U+10FFFF are 2 characters.
	{
		url:      `file:///C:` + strings.Repeat(`/long`, 50) + `/ab.txt`,
		filePath: `C:` + strings.Repeat(`\long`, 50) + `\ab.txt`,
	},
	{
		url:      `file:///C:` + strings.Repeat(`/long`, 50) + `/abX.txt`,
		filePath: `\\?\C:` + strings.Repeat(`\long`, 50) + `\abX.txt`,
	},

	// --- Extended-length UNC (normalize) ---
	// Path longer than MAX_PATH (260 characters including NUL) requires
	// extended-length prefix \\?\UNC.
	// TODO: Characters are UTF-16, U+10000 … U+10FFFF are 2 characters.
	{
		url:      `file://server/share` + strings.Repeat(`/long`, 47) + `/filea.txt`,
		filePath: `\\server\share` + strings.Repeat(`\long`, 47) + `\filea.txt`,
	},
	{
		url:      `file://server/share` + strings.Repeat(`/long`, 47) + `/fileaX.txt`,
		filePath: `\\?\UNC\server\share` + strings.Repeat(`\long`, 47) + `\fileaX.txt`,
	},

	// --- Missing path after scheme ---
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
		url:     `ftp:///C:/Windows/System32`,
		wantErr: `non-file URL`,
	},

	// --- Incorrect disk syntax ---
	{
		url:     `file:///Z|/invalid/path.txt`,
		wantErr: `file URL missing drive letter`,
	},

	// --- Host but no share (UNC malformed) ---
	{
		url:     `file://server/`,
		wantErr: `file URL missing UNC share name`,
	},

	// --- Four slashes without host ---
	{
		url:     `file://///`,
		wantErr: `file URL missing drive letter`,
	},
}
