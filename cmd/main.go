package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/brojonat/temporal-examples/auction"
	"github.com/brojonat/temporal-examples/worker"
	"github.com/urfave/cli/v2"
)

func getDefaultLogger(lvl slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     lvl,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.Function = ""
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	}))
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Run examples.",
				Subcommands: []*cli.Command{
					{
						Name:  "auction-server",
						Usage: "Run the auction server",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "port",
								Aliases:  []string{"p"},
								Required: true,
								Usage:    "Port to listen on",
								Value:    ":8080",
							},
							&cli.StringFlag{
								Name:     "temporal-host",
								Required: true,
								Usage:    "Temporal host",
								Value:    "http://localhost:7223",
							},
						},
						Action: func(ctx *cli.Context) error {
							return auction.RunHTTPServer(
								ctx.Context,
								getDefaultLogger(slog.LevelInfo),
								ctx.String("port"),
								ctx.String("temporal-host"),
							)
						},
					},
					{
						Name:  "worker",
						Usage: "Run the temporal worker",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "temporal-host",
								Required: true,
								Usage:    "Temporal host",
								Value:    "http://localhost:7223",
							},
						},
						Action: func(ctx *cli.Context) error {
							return worker.RunWorker(
								ctx.Context,
								getDefaultLogger(slog.LevelInfo),
								ctx.String("temporal-host"),
							)
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error running command: %s\n", err.Error())
		os.Exit(1)
	}
}
