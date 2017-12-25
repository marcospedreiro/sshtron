package main

import (
	"fmt"
	"os"

	"github.com/marcospedreiro/sshtron/config"
	"github.com/marcospedreiro/sshtron/server"
	"github.com/marcospedreiro/sshtron/version"
)

func main() {

	fmt.Printf("sshtron version: %s\n", version.VERSION)
	cfgFilePath := "config/resources/config.json"
	if len(os.Args) >= 2 {
		cfgFilePath = os.Args[1]
	}

	cfg, err := config.CreateConfig(cfgFilePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		panic("Unable to load configuration")
	}

	//fmt.Printf("%+v\n", cfg)

	server.Start(cfg)

	return
}
