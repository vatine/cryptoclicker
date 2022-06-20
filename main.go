package main

import (
	"flag"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	log "github.com/sirupsen/logrus"

	"github.com/vatine/cryptoclicker/pkg/game"
)

func main() {
	logLevel := flag.String("loglevel", "info", "Log level")
	flag.Parse()

	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("failed to parse error log flag value")
		level = log.WarnLevel
	}
	log.SetLevel(level)

	log.Debug("Starting game now")

	g := game.NewGame()
	fmt.Println(ebiten.RunGame(g))
}
