package http_transport

import (
	"encoding/json"
	"fintech-capstone/m/v2/internal/api_gateway/adapters/transports/http/http_kit/writer"
	"fintech-capstone/m/v2/internal/api_gateway/adapters/transports/http/middleware"
	"fintech-capstone/m/v2/internal/api_gateway/apperr"
	"fintech-capstone/m/v2/internal/api_gateway/gateway"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
	"fintech-capstone/m/v2/internal/api_gateway/types"
	"net/http"
	"net/http/pprof"
)

func NewRouter(gw *gateway.Gateway, logger ports.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /transfer",
		Unary[types.TransferCommand, types.TransferResult](
			gw.TransferHandler,
			TransferJSONDecoder(),
			func(w http.ResponseWriter, res types.TransferResult) {
				writer.JSON(w, http.StatusOK, types.TransferResponse{
					TransactionID: res.TransactionID,
					Status:        res.Status,
					Message:       res.Message,
				})
			},
			DefaultMeta,
		),
	)

	mux.HandleFunc("GET /metrics",
		Unary[struct{}, types.MetricsSnapshot](
			gw.MetricsHandler, // ports.UnaryHandler[struct{}, types.MetricsSnapshot]
			EmptyDecoder,
			JSONEncoder[types.MetricsSnapshot](http.StatusOK),
			DefaultMeta,
		),
	)

	mux.HandleFunc("GET /healthz",
		Unary[struct{}, map[string]string](
			gw.HealthHandler,
			EmptyDecoder,
			JSONEncoder[map[string]string](http.StatusOK),
			DefaultMeta,
		),
	)

	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return middleware.Chain(
		mux,
		middleware.RecovererWithLogger(logger),
		middleware.RequestID,
		middleware.RequestLogger(logger),
		middleware.LimitBytes(1<<20),
	)
}

func TransferJSONDecoder() Decoder[types.TransferCommand] {
	return func(r *http.Request) (types.TransferCommand, error) {
		var dto struct {
			From           string `json:"from"`
			To             string `json:"to"`
			Amount         int64  `json:"amount"`
			IdempotencyKey string `json:"idempotency_key"`
		}
		dec := json.NewDecoder(r.Body) // no MaxBytesReader here
		dec.DisallowUnknownFields()
		if err := dec.Decode(&dto); err != nil {
			return types.TransferCommand{}, apperr.Invalid("invalid JSON payload")
		}
		return types.TransferCommand{
			FromAccount:    dto.From,
			ToAccount:      dto.To,
			AmountCents:    dto.Amount,
			IdempotencyKey: dto.IdempotencyKey,
		}, nil
	}
}
