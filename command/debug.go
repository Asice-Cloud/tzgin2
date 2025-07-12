package command

import (
	"github.com/urfave/cli/v2"
	"github.com/xjtu-tenzor/tz-gin/debugger"
)

func Debug(c *cli.Context) error {
	dbg, err := debugger.NewDebugger("")
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	repl := debugger.NewREPL(dbg)
	repl.Start()

	return nil
}
