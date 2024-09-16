package server

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/poll/temporal"
	"github.com/brojonat/temporal-examples/worker"
	"go.temporal.io/sdk/client"
)

func idFromPrompt(p string) string {
	return fmt.Sprintf("poll: %s", p)
}

// run an http server with endpoints for the auction workflow
func RunHTTPServer(
	ctx context.Context,
	l *slog.Logger,
	port string,
	tcHost string,
) error {

	tc, err := client.Dial(client.Options{
		Logger:   l,
		HostPort: tcHost,
	})
	if err != nil {
		return fmt.Errorf("could not initialize Temporal client: %w", err)
	}
	defer tc.Close()

	mux := http.NewServeMux()
	mux.Handle("POST /start", handleStart(l, tc))
	mux.Handle("POST /vote", handleVote(l, tc))
	mux.Handle("GET /get-state", handleGetState(l, tc))
	mux.Handle("POST /handle-result", handleResult(l, tc))

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

// start a poll
func handleStart(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.RunPollWFRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}

		wopts := client.StartWorkflowOptions{
			ID:        idFromPrompt(payload.Prompt),
			TaskQueue: worker.TaskQueue,
		}
		_, err = tc.ExecuteWorkflow(r.Context(), wopts, temporal.RunPollWF, payload)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		convenience.WriteOK(w)
	}
}

// query the workflow for the current top bid
func handleGetState(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := idFromPrompt(r.URL.Query().Get("prompt"))
		response, err := tc.QueryWorkflow(r.Context(), id, "", temporal.QueryTypeState)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		var result temporal.PollResult
		if err = response.Get(&result); err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		msg := fmt.Sprintf("Poll results for \"%s\":\n", result.Prompt)

		// to iterate over map in order of values, we have to unpack the
		// map into a slice and sort by the votes
		type optVotes struct {
			Option string
			Votes  float64
		}
		ovs := []optVotes{}
		for o, v := range result.Votes {
			ovs = append(ovs, optVotes{Option: o, Votes: v})
		}
		slices.SortFunc(ovs, func(a, b optVotes) int {
			return cmp.Compare(b.Votes, a.Votes)
		})
		for _, ov := range ovs {
			msg += fmt.Sprintf("\t%s: %v\n", ov.Option, ov.Votes)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(convenience.DefaultJSONResponse{Message: msg})
	}
}

// send a signal to the workflow with the supplied bid
func handleVote(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.PollVote
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}

		id := idFromPrompt(payload.Prompt)
		err = tc.SignalWorkflow(r.Context(), id, "", temporal.SignalTypeVote, payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}

		convenience.WriteOK(w)
	}
}

// handle the winning bid webhook
func handleResult(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.PollResult
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}
		l.Info(
			"got poll result",
			"prompt", payload.Prompt,
			"votes", payload.Votes,
		)
		convenience.WriteOK(w)
	}
}
