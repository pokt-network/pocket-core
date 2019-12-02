// +build windows

package config

import "os"

const (
	FILESEPARATOR = "\\"
)

var (
	DATADIR = os.Getenv("APPDATA") + "\\Pocket"
)
