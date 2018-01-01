package main

import (
	"github.com/FreifunkSearchProjekt/Server/common"
)

func main() {
	conf := common.LoadConfigs()
	r := common.Setup()

	common.Begin(r, conf)
}