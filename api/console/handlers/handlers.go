package handlers

import "idiom-api-services/api/console/config"

type Handler struct {
	config config.AppConfig
}

func NewHandler(config config.AppConfig) *Handler {
	return &Handler{
		config: config,
	}
}
