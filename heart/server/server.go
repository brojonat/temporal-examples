package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brojonat/temporal-examples/convenience"
	"github.com/brojonat/temporal-examples/heart/temporal"
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

	listenAddr := fmt.Sprintf(":%s", port)
	l.Info("listening", "port", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

// start a long lived workflow
func handleStart(l *slog.Logger, tc client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wopts := client.StartWorkflowOptions{
			ID:        "heartbeat-and-continue-workflow",
			TaskQueue: worker.TaskQueue,
		}
		_, err := tc.ExecuteWorkflow(r.Context(), wopts, temporal.RunHeartWF)
		if err != nil {
			convenience.WriteInternalError(l, w, err)
			return
		}
		convenience.WriteOK(w)
	}
}
