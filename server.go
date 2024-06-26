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
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

var Version = "dev"

const (
	SHUTDOWN_TIMEOUT    = 5 * time.Second
	READ_HEADER_TIMEOUT = 2 * time.Second
	ERR_EXIT            = 3
)

func main() {
	var listen, level string
	var proxy bool
	flag.StringVar(&level, "verbosity", "info", "log level")
	flag.StringVar(&level, "v", "info", "log level")
	flag.StringVar(&listen, "listen", ":8080", "address to listen on")
	flag.StringVar(&listen, "l", ":8080", "address to listen on")
	flag.BoolVar(&proxy, "proxy", false, "running behind a reverse proxy")

	flag.Parse()
	setupLogging(level)
	server := setupServer(listen, proxy)

	go startServer(server)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), SHUTDOWN_TIMEOUT)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.WithLevel(zerolog.FatalLevel).Err(err).Msg("Error while shutting down server")
		return
	}
}

func setupLogging(userLevel string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if pid := os.Getpid(); pid == 1 {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger() // json text logger when in a container
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if level, err := zerolog.ParseLevel(userLevel); err != nil {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Warn().Err(err).Msg("invalid log level, defaulting to info")
	} else {
		zerolog.SetGlobalLevel(level)
	}
}

func setupServer(addr string, useProxy bool) *http.Server {
	commonMiddleware := []server.Middleware{
		hlog.NewHandler(log.Logger), // just inject the logger in the context here
		hlog.RequestIDHandler("request_id", "Request-Id"),
	}
	if useProxy {
		// proxy must come before RemoteIPHandler, if its being used
		commonMiddleware = append(commonMiddleware, server.ProxyMiddleware)
	}
	commonMiddleware = append(commonMiddleware,
		hlog.RemoteIPHandler("remote_ip"),
		hlog.MethodHandler("method"),
		hlog.HostHandler("host", true),
		hlog.URLHandler("url"),
		hlog.UserAgentHandler("user_agent"),
		server.LoggerMiddleware,
		server.SecurityMiddleware,
		server.HeaderMiddleware,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", server.HealthzHandler)
	mux.HandleFunc("/", server.RedirectHandler) // direct requests to root to /next
	mux.HandleFunc("/next", server.NextHandler)
	return &http.Server{Addr: addr, ReadHeaderTimeout: READ_HEADER_TIMEOUT, Handler: server.WrapHandler(mux, commonMiddleware...), ErrorLog: stdlog.New(log.Logger, "", 0)}
}

func startServer(server *http.Server) {
	// Start the HTTP server on port 8080
	log.Info().Str("addr", server.Addr).Str("version", Version).Msg("starting server")
	// fmt.Println("Server listening on port 8080...")
	log.Fatal().
		Err(server.ListenAndServe()).
		Msg("closing server")
}
