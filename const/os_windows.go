// +build windows

package _const

import "os"

/*
This file is for OS specific data directory configuration
 */
const (
	FILESEPARATOR = "\\"
)

var (
	DATADIR = os.Getenv("APPDATA") + "\\" + "Pocket"
)
