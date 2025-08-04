package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ShamirMaharjan/Student-management-system--GO/internal/config"
	"github.com/ShamirMaharjan/Student-management-system--GO/internal/http/controller/student"
	"github.com/ShamirMaharjan/Student-management-system--GO/storage/sqlite"
)

func main() {

	//load config
	cfg := config.MustLoad()

	//setup database
	_, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Database connected", slog.String("database", cfg.StoragePath))

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.CreateStudent())

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started at port", slog.String("port", cfg.Addr))

	//using channel to synchronize the server like a server stop garyo re tara euta event ko executation complete vako xaina re
	//so tyo event ko executation complete na vaye samma main go routine stop hundaina as channel le block garira hunxa

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	<-done

	//kaile kai server stop hundaina infinitly run vairakhxa tei vayera hamle context use gareko
	//so if server 5 second samma stop vayena vani server stop vayena vanera msg auxa
	//natra server shutdown hunxa
	slog.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server Shutdown Gracefully")
}
