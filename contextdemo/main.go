package main

import (
	_ "contextdemo/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"contextdemo/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
