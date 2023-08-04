//go:generate go run ../mylife-home-core-generator/cmd/main.go .

package main

import (
	_ "mylife-home-plugins-logic-base/plugin"
)

// Makefile :
// go generate
// go build -buildmode=plugin -o logic-base.so main.go
