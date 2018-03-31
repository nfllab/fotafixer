fotafixer
=========

This tool can fix obfuscated Samsung Android firmware OTA updates so that the
update zip file can be opened successfully.

Build
-----

`go build fotafixer.go` (or `go get github.com/nfllab/fotafixer`)

Usage
-----

Usage: ./fotafixer inputfile [outputfile]

If you supply an outputfile parameter, then the inputfile is unchanged.
Otherwise the inputfile is fixed in place.
