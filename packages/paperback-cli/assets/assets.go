package assets

import _ "embed"

//go:embed homepage.html
var HomepageHTML []byte

//go:embed shims.js
var ShimsJS []byte
