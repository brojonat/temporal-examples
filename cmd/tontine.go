package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/urfave/cli/v2"
	"github.com/brojonat/temporal-examples/tontine/server"
)

func tontine_run_server(ctx *cli.Context) error {
	return server.RunHTTPServer(
		ctx.Context,
		getDefaultLogger(slog.LevelInfo),
		ctx.String("port"),
		ctx.String("temporal-host"),
	)
}

func start_tontine(ctx *cli.Context) error {
	body := struct {
		Name         string   `json:"name"`
		Participants []string `json:"participants"`
	}{
		Name:         ctx.String("name"),
		Participants: ctx.StringSlice("participant"),
	}

	if len(body.Name) < 1 {
		return fmt.Errorf("must supply a tontine name")
	}
	if len(body.Participants) < 2 {
		return fmt.Errorf("must supply at least two participants")
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
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response code (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
}

func signal_participant_alive(ctx *cli.Context) error {
	body := struct {
		Name       string `json:"name"`
		Participant string `json:"participant"`
		Alive      bool   `json:"alive"`
	}{
		Name:       ctx.String("name"),
		Participant: ctx.String("participant"),
		Alive:      ctx.Bool("alive"),
	}

	if len(body.Name) < 1 {
		return fmt.Errorf("must supply a tontine name")
	}
	if len(body.Participant) < 1 {
		return fmt.Errorf("must supply a participant name")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, ctx.String("endpoint")+"/signal-alive", bytes.NewReader(b))
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
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("bad response code (%d) and error reading body: %w", res.StatusCode, err)
	}
	return fmt.Errorf("bad response code (%d): %s", res.StatusCode, b)
}

func get_tontine_state(ctx *cli.Context) error {
	r, err := http.NewRequest(http.MethodGet, ctx.String("endpoint")+"/get-state", nil)
	if err != nil {
		return err
	}
	q := r.URL.Query()
	q.Add("name", ctx.String("name"))
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
	var state struct {
		Participants []string `json:"participants"`
		TotalFund    float64  `json:"totalFund"`
	}
	if err := json.Unmarshal(b, &state); err != nil {
		return fmt.Errorf("could not parse message: %w: %s", err, b)
	}
	fmt.Printf("Tontine State:\nParticipants: %v\nTotal Fund: %.2f\n", state.Participants, state.TotalFund)
	return nil
}
