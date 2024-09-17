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
							return auction_run_server(ctx)
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
							return auction_run_worker(ctx)
						},
					},
					{
						Name:  "start",
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
								Value:   "http://localhost:8080/webhook",
							},
						},
						Action: func(ctx *cli.Context) error {
							return start_auction(ctx)
						},
					},
					{
						Name:  "get-state",
						Usage: "Get the current state of the auction",
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
							return get_auction_state(ctx)
						},
					},
					{
						Name:  "bid",
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
							return auction_place_bid(ctx)
						},
					},
				},
			},
			{
				Name:  "poll",
				Usage: "poll related subcommands",
				Subcommands: []*cli.Command{
					{
						Name:  "run-server",
						Usage: "Run the poll server",
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
							return poll_run_server(ctx)
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
							return poll_run_worker(ctx)
						},
					},
					{
						Name:  "start",
						Usage: "start a poll",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP server endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "prompt",
								Required: true,
								Aliases:  []string{"p", "q"},
								Usage:    "Prompt question for the poll",
							},
							&cli.StringSliceFlag{
								Name:     "option",
								Required: true,
								Aliases:  []string{"opt", "o"},
								Usage:    "Options for poll",
							},
							&cli.StringFlag{
								Name:     "duration",
								Required: true,
								Aliases:  []string{"dur", "d"},
								Usage:    "Poll duration in Go time.Duration format (e.g., 15m)",
							},
							&cli.StringFlag{
								Name:    "webhook",
								Aliases: []string{"web", "w"},
								Usage:   "Webhook endpoint for poll results",
								Value:   "http://localhost:8080/webhook",
							},
						},
						Action: func(ctx *cli.Context) error {
							return start_poll(ctx)
						},
					},
					{
						Name:  "get-state",
						Usage: "Get the current state of the poll",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "prompt",
								Required: true,
								Aliases:  []string{"p"},
								Usage:    "Prompt of the poll to fetch results for",
							},
						},
						Action: func(ctx *cli.Context) error {
							return get_poll_state(ctx)
						},
					},
					{
						Name:  "vote",
						Usage: "vote on a poll",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "prompt",
								Required: true,
								Aliases:  []string{"p"},
								Usage:    "Prompt of the poll to vote on",
							},
							&cli.StringFlag{
								Name:     "option",
								Required: true,
								Aliases:  []string{"opt", "o"},
								Usage:    "Option to cast vote for",
							},
							&cli.Float64Flag{
								Name:     "amount",
								Required: true,
								Aliases:  []string{"a"},
								Usage:    "Magnitude of vote",
								Value:    1,
							},
						},
						Action: func(ctx *cli.Context) error {
							return poll_vote(ctx)
						},
					},
				},
			},
			{
				Name:  "dms",
				Usage: "DMS related subcommands",
				Subcommands: []*cli.Command{
					{
						Name:  "run-server",
						Usage: "Run the DMS server",
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
							return dms_run_server(ctx)
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
							return dms_run_worker(ctx)
						},
					},
					{
						Name:  "start",
						Usage: "start a DMS",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP server endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "id",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "ID for the DMS",
							},
							&cli.StringFlag{
								Name:     "message",
								Required: true,
								Aliases:  []string{"m"},
								Usage:    "Message to send as contingency",
							},
							&cli.StringFlag{
								Name:     "duration",
								Required: true,
								Aliases:  []string{"dur", "d"},
								Usage:    "DMS duration in Go time.Duration format (e.g., 15m)",
							},
							&cli.StringFlag{
								Name:    "webhook",
								Aliases: []string{"web", "w"},
								Usage:   "Webhook endpoint for contingency message",
								Value:   "http://localhost:8080/webhook",
							},
						},
						Action: func(ctx *cli.Context) error {
							return start_dms(ctx)
						},
					},
					{
						Name:  "get-state",
						Usage: "Get the current state of the DMS",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "id",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "ID for the DMS",
							},
						},
						Action: func(ctx *cli.Context) error {
							return get_dms_state(ctx)
						},
					},
					{
						Name:  "deactivate",
						Usage: "deactivate a DMS",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "endpoint",
								Usage: "HTTP endpoint",
								Value: "http://localhost:8080",
							},
							&cli.StringFlag{
								Name:     "id",
								Required: true,
								Aliases:  []string{"i"},
								Usage:    "ID for the DMS",
							},
						},
						Action: func(ctx *cli.Context) error {
							return dms_deactivate(ctx)
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
