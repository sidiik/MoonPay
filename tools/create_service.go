package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

var folders = []string{
	"cmd/",
	"internal/domain",
	"internal/usecase",
	"internal/repository",
	"internal/delivery/http",
	"internal/delivery/grpc",
	"pkg",
}

func main() {
	slog.Info("Creating service")
	if len(os.Args) < 2 {
		log.Fatalf("Usage: go run create_service.go <service_name>")
	}

	serviceName := strings.ToLower(os.Args[1])
	createServiceArchitecture(serviceName)

}

func createServiceArchitecture(serviceName string) {
	for _, folder := range folders {
		path := fmt.Sprintf("services/%s/%s", serviceName, folder)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("Failed to create folder %s: %v", path, err)
		}
	}

	fmt.Printf("✅ Clean Architecture structure for service '%s' created successfully.\n", serviceName)
}
