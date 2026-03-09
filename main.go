package main

import (
	"file_graph/internal/handlers"
	"file_graph/internal/logger"
	"fmt"
	"log"
	"net/http"
)

func main() {
	logger.Println("--- START ---")

	http.HandleFunc("/", handlers.ServeIndex)
	http.HandleFunc("/favicon.ico", handlers.ServeFavicon)
	http.HandleFunc("/api/scan", handlers.HandleScan)
	http.HandleFunc("/api/open", handlers.HandleOpen)
	http.HandleFunc("/api/log", handlers.HandleClientLog)
	http.HandleFunc("/api/heartbeat", handlers.HandleHeartbeat)
	http.HandleFunc("/api/progress", handlers.HandleProgress)
	http.HandleFunc("/api/shutdown", handlers.HandleShutdown)
	http.HandleFunc("/api/kill", handlers.HandleHardExit)

	port := "8080"
	fmt.Printf("Starting server on http://localhost:%s\n", port)
	logger.Printf("Server: http://localhost:%s\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
