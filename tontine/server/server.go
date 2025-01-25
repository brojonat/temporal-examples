package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/tontine/temporal"
	"go.temporal.io/sdk/client"
)

// idFromTontineName generates a unique workflow ID for a tontine based on its name.
func idFromTontineName(name string) string {
	return fmt.Sprintf("tontine:%s", name)
}

// RunHTTPServer runs an HTTP server with endpoints to manage the tontine workflow.
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
	mux.Handle("/start", handleStart(l, tc))
	mux.Handle("/get-state", handleGetState(l, tc))
	mux.Handle("/signal-alive", handleSignalAlive(l, tc))

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

// handleStart initializes a new tontine workflow.
func handleStart(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Name         string   `json:"name"`
			Participants []string `json:"participants"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}

		wopts := client.StartWorkflowOptions{
			ID:        idFromTontineName(payload.Name),
			TaskQueue: "tontine_task_queue",
		}
		_, err := tc.ExecuteWorkflow(r.Context(), wopts, temporal.TontineWorkflow, payload.Participants)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		convenience.WriteOK(w)
	}
}

// handleGetState queries the current state of the tontine workflow.
func handleGetState(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		id := idFromTontineName(name)
		response, err := tc.QueryWorkflow(r.Context(), id, "", temporal.QueryTontineState)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		var state struct {
			Participants []string `json:"participants"`
			TotalFund    float64  `json:"totalFund"`
		}
		if err := response.Get(&state); err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(state)
	}
}

// handleSignalAlive sends a signal to update a participant's alive status.
func handleSignalAlive(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Name        string `json:"name"`
			Participant string `json:"participant"`
			Alive       bool   `json:"alive"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}

		id := idFromTontineName(payload.Name)
		err := tc.SignalWorkflow(r.Context(), id, "", temporal.SignalParticipantAlive, payload)
		if err != nil {
			convenience.WriteBadRequestError(w, err)
			return
		}
		convenience.WriteOK(w)
	}
}
