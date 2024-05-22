package main

import (
	"context"
	"flag"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aztechian/nextride-shortcut/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var listen, level string
	flag.StringVar(&level, "verbosity", "info", "log level")
	flag.StringVar(&level, "v", "info", "log level")
	flag.StringVar(&listen, "listen", ":8080", "address to listen on")
	flag.StringVar(&listen, "l", ":8080", "address to listen on")

	flag.Parse()
	setupLogging(level)
	server := setupServer(listen)

	go startServer(server)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error while shutting down server")
	}
}

func setupLogging(userLevel string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if level, err := zerolog.ParseLevel(userLevel); err != nil {
		zerolog.SetGlobalLevel(level)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupServer(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.RedirectHandler) // direct requests to root to /next
	mux.HandleFunc("/next", server.NextHandler)
	return &http.Server{Addr: addr, Handler: server.Logger(mux), ErrorLog: stdlog.New(log.Logger, "", 0)}
}

func startServer(server *http.Server) {
	// Start the HTTP server on port 8080
	log.Info().Str("addr", server.Addr).Msg("starting server")
	// fmt.Println("Server listening on port 8080...")
	log.Fatal().
		Err(server.ListenAndServe()).
		Msg("closing server")
}
