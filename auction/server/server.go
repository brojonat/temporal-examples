package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/temporal-examples/auction/temporal"
	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/worker"
	"go.temporal.io/sdk/client"
)

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
	mux.Handle("POST /bid", handleBid(l, tc))
	mux.Handle("GET /get-state", handleGetState(l, tc))
	mux.Handle("POST /webhook", handleResult(l, tc))

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

// start an auction
func handleStart(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.RunAuctionWFRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		wopts := client.StartWorkflowOptions{
			ID:        payload.Item,
			TaskQueue: worker.TaskQueue,
		}
		_, err = tc.ExecuteWorkflow(r.Context(), wopts, temporal.RunAuctionWF, payload)
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

		response, err := tc.QueryWorkflow(
			r.Context(), r.URL.Query().Get("item"), "", temporal.QueryTypeState)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		var result temporal.QueryResultState
		if err = response.Get(&result); err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		msg := fmt.Sprintf("top bid by %s for %f", result.Bidder, result.Amount)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(convenience.DefaultJSONResponse{Message: msg})
	}
}

// send a signal to the workflow with the supplied bid
func handleBid(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.AuctionBid
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}

		err = tc.SignalWorkflow(r.Context(), payload.Item, "", temporal.SignalTypeBid, payload)
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

		var payload temporal.AuctionBid
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}

		l.Info(
			"got auction result",
			"item", payload.Item,
			"bidder", payload.Bidder,
			"amount", payload.Amount,
		)
		convenience.WriteOK(w)
	}
}
