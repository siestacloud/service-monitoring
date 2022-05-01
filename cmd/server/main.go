package main

import (
	"fmt"
	"log"
	"os"

	"github.com/siestacloud/service-monitoring/internal/server"
	"github.com/siestacloud/service-monitoring/internal/server/config"
)

var (
	cli config.CLI
	cfg config.ServerConfig
)

func main() {
	fmt.Println("HERE")
	err := config.Parse(cli.ConfigPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg)
	s, err := server.New(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = s.Start(); err != nil {
		os.Exit(0)
	}
	os.Exit(0)
}
