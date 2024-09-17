package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/dms/temporal"
	"github.com/brojonat/temporal-examples/worker"
	"go.temporal.io/sdk/client"
)

func idFromID(id string) string {
	return fmt.Sprintf("dms: %s", id)
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
	mux.Handle("POST /deactivate", handleDeactivate(l, tc))
	mux.Handle("GET /get-state", handleGetState(l, tc))
	mux.Handle("POST /webhook", handleResult(l, tc))

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

// start a dms
func handleStart(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.RunDMSWFRequest
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		if payload.ID == "" || payload.Message == "" {
			convenience.WriteBadRequestError(w, fmt.Errorf("must supply dms id and message"))
			return
		}
		wopts := client.StartWorkflowOptions{
			ID:        idFromID(payload.ID),
			TaskQueue: worker.TaskQueue,
		}
		_, err = tc.ExecuteWorkflow(r.Context(), wopts, temporal.RunDMSWF, payload)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		convenience.WriteOK(w)
	}
}

// query the dms for the current state
func handleGetState(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := idFromID(r.URL.Query().Get("id"))
		response, err := tc.QueryWorkflow(r.Context(), id, "", temporal.QueryTypeState)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		var result string
		if err = response.Get(&result); err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(convenience.DefaultJSONResponse{Message: result})
	}
}

// deactivate the dms
func handleDeactivate(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := idFromID(r.URL.Query().Get("id"))
		err := tc.SignalWorkflow(r.Context(), id, "", temporal.SignalTypeDeactivate, nil)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}
		convenience.WriteOK(w)
	}
}

// handle the dms timeout
func handleResult(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload temporal.DMSTimeoutPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}
		l.Info(
			"got dms timeout",
			"message", payload.Message,
		)
		convenience.WriteOK(w)
	}
}
