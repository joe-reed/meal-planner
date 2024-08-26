//go:build tools

// The below are used at build time. We include them here to
// ensure that they are not removed by go mod tidy

package main

import (
	_ "github.com/campoy/jsonenums"
)
