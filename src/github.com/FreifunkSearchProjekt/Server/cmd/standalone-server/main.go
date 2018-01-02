package main

import (
	"github.com/FreifunkSearchProjekt/Server/common"
	"log"
)

func main() {
	conf := common.LoadConfigs()
	r := common.Setup()

	log.Println("Starting SearchServer")
	common.Begin(r, conf)
}
