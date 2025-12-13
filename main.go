package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alireza12prom/intellimerge/internal/config"
	"github.com/alireza12prom/intellimerge/internal/webhook"
	"github.com/caarlos0/env/v11"
)

func main() {
	cfg := &config.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	handler := webhook.NewHandler(cfg)

	http.HandleFunc(cfg.WebhookPath, handler.HandleWebhook)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting IntelliMerge server on %s", addr)
	log.Printf("Webhook endpoint: %s", cfg.WebhookPath)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
