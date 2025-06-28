// main.go
package main

import (
	"github.com/alecthomas/kong"
	"github.com/williampsena/ebenezer-cli/cmd"
	internalcmd "github.com/williampsena/ebenezer-cli/internal/cmd"
)

func main() {
	cli := cmd.CLI{}
	ctx := kong.Parse(&cli)
	err := ctx.Run(&internalcmd.Context{
		Debug: cli.Debug,
	})
	ctx.FatalIfErrorf(err)
}
