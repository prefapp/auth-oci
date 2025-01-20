package main

import (
	"github.com/alecthomas/kong"
)

func main() {

	kongCtx := kong.Parse(&CLI)

	switch kongCtx.Command() {

	case "login":
		multilogin()

	case "validate":
		validate()

	default:
		panic(kongCtx.Command())
	}

}
