package entrypoint

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports/inbound"
)

func (g *Gateway) HealthHandler(context.Context, inbound.RequestMeta, struct{}) (map[string]string, error) {
	return map[string]string{"status": "ok"}, nil
}
