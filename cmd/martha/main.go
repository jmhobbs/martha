package main

import (
	"flag"

	"github.com/jmhobbs/martha/internal/app"
)

func main() {
	var (
		configPath *string = flag.String("config", "martha.yaml", "Configuration file, should be writable by this process.")
		debug *bool = flag.Bool("debug", true, "sets log level to debug")
	)
	flag.Parse()

	app.Run(*debug, *configPath)
}
