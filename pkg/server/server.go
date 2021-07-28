package server

import (
	"fmt"
	"net/http"

	config2 "github.com/oakeshq/go-starter/config"
	"go.uber.org/zap"
)

// ListenAndServe starts the REST API.
func ListenAndServe(cfg *config2.Config, h http.Handler) error {
	location := fmt.Sprintf("%v:%v", cfg.SERVER.Host, cfg.SERVER.Port)
	zap.L().Sugar().Infof("API started on: %s", location)
	return http.ListenAndServe(location, h)
}

