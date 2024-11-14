package http_handlers

import (
	"log/slog"

	"github.com/rozhnof/stakewolle-auth-service/internal/application/services"
	"go.opentelemetry.io/otel/trace"
)

type AuthHandler struct {
	log         *slog.Logger
	userService *services.UserService
	tracer      trace.Tracer
}

func NewAuthHandler(service *services.UserService, log *slog.Logger, tracer trace.Tracer) *AuthHandler {
	return &AuthHandler{
		userService: service,
		log:         log,
		tracer:      tracer,
	}
}
