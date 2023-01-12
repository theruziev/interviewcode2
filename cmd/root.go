package cmd

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type Ctx struct {
	IsDebug  bool
	LogLevel string
}

var Root struct {
	IsDebug  bool   `help:"Enable Debug mode" env:"DEBUG"`
	LogLevel string `help:"Log level" default:"debug" env:"LOG_LEVEL"`

	Httpserver httpserver `cmd:""`
	UserEmail  userEmail  `cmd:""`
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to read env: %s", err)
	}
	rootCmd := &Root
	ctx := kong.Parse(rootCmd)

	ctx.FatalIfErrorf(ctx.Run(&Ctx{
		IsDebug:  rootCmd.IsDebug,
		LogLevel: rootCmd.LogLevel,
	}))
}
