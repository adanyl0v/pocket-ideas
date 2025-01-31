package app

import (
	"github.com/adanyl0v/pocket-ideas/internal/config"
	"log"
)

func Run() {
	cfg := config.MustReadFile(config.DefaultFilePath())
	log.Printf("read config %+v\n", *cfg)
}
