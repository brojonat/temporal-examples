package convenience

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime"
	"time"
)

type DefaultJSONResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func WriteOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(DefaultJSONResponse{Message: "ok"})
}

func WriteBadRequestError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	resp := DefaultJSONResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(resp)
}

func WriteEmptyResultError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	resp := DefaultJSONResponse{Error: "empty result set"}
	json.NewEncoder(w).Encode(resp)
}

func WriteInternalError(l *slog.Logger, w http.ResponseWriter, e error) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelError, e.Error(), pcs[0])
	_ = l.Handler().Handle(context.Background(), r)
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(DefaultJSONResponse{Error: "internal error"})
}
