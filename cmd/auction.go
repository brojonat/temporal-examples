package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/temporal-examples/auction/server"
	"github.com/brojonat/temporal-examples/auction/temporal"
	"github.com/brojonat/temporal-examples/worker"
	"github.com/urfave/cli/v2"
)

func run_server(ctx *cli.Context) error {
	return server.RunHTTPServer(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("port"),
		ctx.String("temporal-host"),
	)
}

func run_worker(ctx *cli.Context) error {
	return worker.RunWorker(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("temporal-host"),
	)
}

func start_auction(ctx *cli.Context) error {
	dur, err := time.ParseDuration(ctx.String("duration"))
	if err != nil {
		return err
	}
	body := temporal.RunAuctionWFRequest{
		StartTime:    time.Now(),
		Duration:     dur,
		Item:         ctx.String("item"),
		ReservePrice: ctx.Float64("reserve-price"),
		Webhook:      ctx.String("webhook"),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/start-auction", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response code (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
}

func place_bid(ctx *cli.Context) error {
	body := temporal.AuctionBid{
		Item:   ctx.String("item"),
		Bidder: ctx.String("bidder"),
		Amount: ctx.Float64("amount"),
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/place-bid", bytes.NewReader(b))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusOK {
		return nil
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response code (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
}

func get_top_bid(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/get-top-bid", nil)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	q.Add("item", ctx.String("item"))
	r.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading body: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
	}
	fmt.Printf("%s\n", b)
	return nil
}
