package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/poll/server"
	"github.com/brojonat/temporal-examples/poll/temporal"
	"github.com/brojonat/temporal-examples/worker"
	"github.com/urfave/cli/v2"
)

func poll_run_server(ctx *cli.Context) error {
	return server.RunHTTPServer(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("port"),
		ctx.String("temporal-host"),
	)
}

func poll_run_worker(ctx *cli.Context) error {
	return worker.RunWorker(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("temporal-host"),
	)
}

func start_poll(ctx *cli.Context) error {
	dur, err := time.ParseDuration(ctx.String("duration"))
	if err != nil {
		return err
	}
	body := temporal.RunPollWFRequest{
		StartTime: time.Now(),
		Duration:  dur,
		Prompt:    ctx.String("prompt"),
		Options:   ctx.StringSlice("option"),
		Webhook:   ctx.String("webhook"),
	}
	if len(body.Prompt) < 1 {
		return fmt.Errorf("must supply a poll prompt")
	}
	if len(body.Options) < 1 {
		return fmt.Errorf("must supply at least one option")
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/start", bytes.NewReader(b))
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

func poll_vote(ctx *cli.Context) error {
	body := temporal.PollVote{
		Prompt: ctx.String("prompt"),
		Option: ctx.String("option"),
		Amount: ctx.Float64("amount"),
	}
	if len(body.Prompt) < 1 {
		return fmt.Errorf("must supply a poll prompt")
	}
	if body.Option == "" {
		return fmt.Errorf("must specify option")
	}
	if body.Amount < 0 {
		return fmt.Errorf("cannot vote a negative amount")
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/vote", bytes.NewReader(b))
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

func get_poll_state(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/get-state", nil)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	q.Add("prompt", ctx.String("prompt"))
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
	var body convenience.DefaultJSONResponse
	err = json.Unmarshal(b, &body)
	if err != nil {
		return fmt.Errorf("could not parse message: %w: %s", err, b)
	}
	fmt.Printf(body.Message)
	return nil
}
