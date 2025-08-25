package command

import (
	"github.com/Asice-Cloud/tzgin2/debugger"
	"github.com/urfave/cli/v2"
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
