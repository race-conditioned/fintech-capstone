package http_transport

import (
	"encoding/json"
	"errors"
	"fintech-capstone/m/v2/internal/api_gateway/contracts"
	"fintech-capstone/m/v2/internal/api_gateway/entrypoint"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
	"fintech-capstone/m/v2/internal/platform"
	"fintech-capstone/m/v2/internal/platform/apperr"
	"fintech-capstone/m/v2/internal/platform/http_kit/limiter"
	"fintech-capstone/m/v2/internal/platform/http_kit/middleware"
	"fintech-capstone/m/v2/internal/platform/http_kit/writer"
	"net/http"
	"net/http/pprof"
)

// NewRouter creates and returns a new HTTP router for the API Gateway.
func NewRouter(gw *entrypoint.Gateway, logger platform.Logger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /transfer",
		Unary[inbound.TransferCommand, inbound.TransferResult](
			gw.TransferHandler,
			TransferJSONDecoder(),
			func(w http.ResponseWriter, res inbound.TransferResult) {
				writer.JSON(w, http.StatusOK, contracts.TransferResponse{
					TransactionID: res.TransactionID().String(),
					Status:        res.Status().String(),
					Message:       res.Message(),
				})
			},
			DefaultMeta,
		),
	)

	mux.HandleFunc("GET /metrics",
		Unary[struct{}, contracts.MetricsSnapshot](
			gw.MetricsHandler, // ports.UnaryHandler[struct{}, types.MetricsSnapshot]
			EmptyDecoder,
			JSONEncoder[contracts.MetricsSnapshot](http.StatusOK),
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

	lightLimiter := limiter.NewLightLimiter(10, 20)

	return middleware.Chain(
		mux,
		middleware.RecovererWithLogger(logger),
		middleware.RequestID,
		middleware.RequestLogger(logger),
		middleware.RateLimitHTTP(lightLimiter.Allow),
		middleware.MaxInFlight(1024),
		middleware.LimitBytes(1<<20),
	)
}

// TransferJSONDecoder decodes a TransferCommand from a JSON HTTP request.
func TransferJSONDecoder() Decoder[inbound.TransferCommand] {
	return func(r *http.Request) (inbound.TransferCommand, error) {
		var dto struct {
			From           string `json:"from"`
			To             string `json:"to"`
			Amount         int64  `json:"amount"`
			IdempotencyKey string `json:"idempotency_key"`
		}
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&dto); err != nil {
			var maxErr *http.MaxBytesError
			if errors.As(err, &maxErr) {
				return inbound.TransferCommand{}, apperr.PayloadTooLarge("request body too large")
			}
			return inbound.TransferCommand{}, apperr.Invalid("invalid JSON payload")
		}
		return inbound.NewTransferCommand(
			dto.From,
			dto.To,
			dto.Amount,
			dto.IdempotencyKey,
		), nil
	}
}
