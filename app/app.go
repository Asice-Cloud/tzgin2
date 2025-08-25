package app

import (
	"os"

	"github.com/Asice-Cloud/tzgin2/command"
	"github.com/Asice-Cloud/tzgin2/config"
	"github.com/Asice-Cloud/tzgin2/util"

	"github.com/urfave/cli/v2"

	_ "embed"
)

func InitApp(configString string) *cli.App {
	cfg := config.New()

	app := cli.NewApp()
	cfg.Load(app, configString)
	cfgStruct := cfg.Parse(configString)

	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create operations",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "directory",
					Aliases: []string{"d"},
					Usage:   "the directory you want to generate",
				},
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"r"},
					Usage:   "the remote address you want to pull",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.String("remote") == "" {
					if err := ctx.Set("remote", cfgStruct.RemoteAddr); err != nil {
						return err
					}
				}
				return command.Create(ctx)
			},
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update operations",
			Action:  command.Update,
		},
		{
			Name:    "run",
			Aliases: []string{"r", "ru"},
			Usage:   "run tz-gin project",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "directory",
					Aliases: []string{"d"},
					Usage:   "specify directory directory",
				},
			},
			Action: func(ctx *cli.Context) error {
				if ctx.String("directory") == "" {
					if err := ctx.Set("directory", "."); err != nil {
						return err
					}
				}
				return command.Run(ctx)
			},
		},
		{
			Name:    "debug",
			Aliases: []string{"db"},
			Usage:   "interactive debugger (like gdb/dlv)",
			Action:  command.Debug,
		},
	}

	app.ExitErrHandler = func(cCtx *cli.Context, err error) {
		if err != nil {
			if exitErr, ok := err.(cli.ExitCoder); ok {
				switch exitErr.ExitCode() {
				case 1:
					util.ErrMsg(err.Error() + "\n")
					os.Exit(1)
				case 2:
					util.WarnMsg(err.Error() + "\n") //健康的错误
					os.Exit(0)
				}
			}
		}
	}

	app.CommandNotFound = func(cCtx *cli.Context, str string) {
		util.ErrMsg("command not found: " + str + " can not be recognized as a valid command\n")
		os.Exit(0)
	}

	return app
}
