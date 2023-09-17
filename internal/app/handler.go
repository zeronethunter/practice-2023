package app

import "github.com/labstack/echo/v4"

type Handler struct {
	cfg Config
}

func NewHandler(cfg Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

func (h *Handler) HandlerWebDAV(c echo.Context) error {
	h.cfg.Handler.ServeHTTP(c.Response().Writer, c.Request())

	return nil
}
