package main

import (
	_ "embed"
	"os"

	"github.com/Asice-Cloud/tzgin2/app"
	"github.com/Asice-Cloud/tzgin2/util"
)

//go:embed config.toml
var configStirng string

//go:embed banner.txt
var banner string

func main() {
	app := app.InitApp(configStirng)
	util.SuccessMsg(banner)

	app.Run(os.Args)

}
