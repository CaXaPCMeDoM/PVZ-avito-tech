package main

import (
	"PVZ-avito-tech/config"
	"PVZ-avito-tech/internal/app"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
