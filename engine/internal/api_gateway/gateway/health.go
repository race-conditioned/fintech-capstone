package gateway

import (
	"context"
	"fintech-capstone/m/v2/internal/api_gateway/ports"
)

func (g *Gateway) HealthHandler(context.Context, ports.RequestMeta, struct{}) (map[string]string, error) {
	return map[string]string{"status": "ok"}, nil
}
