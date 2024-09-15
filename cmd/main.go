package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
				Name:  "auction",
				Usage: "Auction related commands.",
				Subcommands: []*cli.Command{
					{
						Name:  "run-server",
						Usage: "Run the auction server",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "port",
								Aliases: []string{"p"},
								Usage:   "Port to listen on",
								Value:   "8080",
							},
							&cli.StringFlag{
								Name:  "temporal-host",
								Usage: "Temporal host",
								Value: "localhost:7233",
							},
						},
						Action: func(ctx *cli.Context) error {
							return run_server(ctx)
						},
					},
					{
						Name:  "run-worker",
						Usage: "Run the temporal worker",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "temporal-host",
								Usage: "Temporal host",
								Value: "localhost:7233",
							},
						},
						Action: func(ctx *cli.Context) error {
							return run_worker(ctx)
						},
					},
					{
						Name:  "start-auction",
						Usage: "start an auction",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP server endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "item",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "Item to auction",
							},
							&cli.Float64Flag{
								Name:     "reserve-price",
								Required: true,
								Aliases:  []string{"reserve", "r"},
								Usage:    "Reserve price of the auction",
							},
							&cli.StringFlag{
								Name:     "duration",
								Required: true,
								Aliases:  []string{"dur", "d"},
								Usage:    "Auction duration in Go time.Duration format (e.g., 15m)",
							},
							&cli.StringFlag{
								Name:    "webhook",
								Aliases: []string{"web", "w"},
								Usage:   "Webhook endpoint for auction results",
								Value:   "http://localhost:8080/handle-winner-bid",
							},
						},
						Action: func(ctx *cli.Context) error {
							return start_auction(ctx)
						},
					},
					{
						Name:  "get-top-bid",
						Usage: "Get the top bid in an auction",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "item",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "Item to auction",
							},
						},
						Action: func(ctx *cli.Context) error {
							return get_top_bid(ctx)
						},
					},
					{
						Name:  "place-bid",
						Usage: "place a bid in an auction",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "item",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "Item to auction",
							},
							&cli.StringFlag{
								Name:     "bidder",
								Required: true,
								Aliases:  []string{"b"},
								Usage:    "Email for the bid",
							},
							&cli.Float64Flag{
								Name:     "amount",
								Required: true,
								Aliases:  []string{"a"},
								Usage:    "Amount to bid",
							},
						},
						Action: func(ctx *cli.Context) error {
							return place_bid(ctx)
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
