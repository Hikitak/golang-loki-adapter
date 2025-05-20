package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-loki-adapter.local/internal/config"
	"golang-loki-adapter.local/internal/database"
	"golang-loki-adapter.local/internal/loki"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbManager, err := database.NewDBManager(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbManager.Close()

	lokiClient := loki.NewLokiClient(&cfg.Loki)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			records, err := dbManager.ProcessQueue()
			if err != nil {
				log.Printf("Error processing queue: %v", err)
				continue
			}

			if len(records) == 0 {
				continue
			}

			// После успешной отправки в Loki
			if err := lokiClient.Send(records); err != nil {
				log.Printf("Failed to send to Loki: %v", err)
				// Обновляем попытки уже выполнено в ProcessQueue
				continue
			}

			// Если отправка успешна - удаляем записи
			if err := dbManager.DeleteProcessed(records); err != nil {
				log.Printf("Failed to delete processed records: %v", err)
			}
		case sig := <-sigChan:
			log.Printf("Received signal: %s. Shutting down", sig)
			return
		}
	}
}