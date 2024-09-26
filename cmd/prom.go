package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/brojonat/temporal-examples/prom/server"
	// NOTE: worker package differs from other examples!
	"github.com/brojonat/temporal-examples/prom/worker"
	"github.com/urfave/cli/v2"
)

func prom_run_server(ctx *cli.Context) error {
	return server.RunHTTPServer(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("port"),
		ctx.String("temporal-host"),
	)
}

func prom_run_worker(ctx *cli.Context) error {
	return worker.RunWorker(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("temporal-host"),
	)
}

func start_prom(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/start", nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response code (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
}
