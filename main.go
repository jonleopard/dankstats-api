package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/nicklaw5/helix"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffyaml"
)

const version = "1.0.0"

type config struct {
	listenAddr int
}

// application sets up our API. Could rename this to 'server', but application
// works too
type application struct {
	config config
	logger *log.Logger
}

func main() {
	// TODO: Implement Zap logger
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	// ========================================================== Flags
	fs := flag.NewFlagSet("DANKSTATS_API", flag.ExitOnError)

	// TODO: Implement config struct
	// var cfg config

	var (
		listenAddr     = fs.String("listen-addr", "localhost:4000", "sets listen address")
		clientID       = fs.String("client-id", "", "sets twitch client ID")
		clientSecret   = fs.String("client-secret", "", "sets twitch client secret")
		appAccessToken = fs.String("app-access-token", "", "sets twitch app access token")
		_              = fs.String("config", "dankstats_api.prod.yaml", "sets config file")
		//dsn        = fs.String("dsn", "", "Postgres DSN")
	)

	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("DANKSTATS_API"),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ffyaml.Parser),
	)
	fmt.Println(*clientID, *clientSecret, *appAccessToken)

	// ========================================================== API SETUP

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	mux := chi.NewRouter()

	app := &application{
		//cfg:    cfg,
		logger: logger,
	}

	srv := &http.Server{
		Handler:      mux,
		Addr:         *listenAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Setup middleware/CORS
	mux.Use(middleware.RequestID)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(render.SetContentType(render.ContentTypeJSON))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// ========================================================== ROUTES SETUP
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/get-token", app.HandleAppToken)
	mux.HandleFunc("/top-games", app.HandleTopGames)
	mux.HandleFunc("/top-channels", app.HandleTopChannels)

	// ========================================================== BOOT
	logger.Printf("API listening on %s", *listenAddr)
	err := srv.ListenAndServe()
	logger.Fatal(err)

	return nil
}

// HandleAppToken retrieves a twitch AppAccessToken for use in subsequent API
// calls.
func (a *application) HandleAppToken(w http.ResponseWriter, r *http.Request) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     "...",
		ClientSecret: "...",
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.RequestAppAccessToken([]string{"analytics:read:games"})
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(resp)
	fmt.Printf("%+v\n", resp)

	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp.GetRateLimitReset())
}

// HandleTopGames responds with the top twitch games
func (a *application) HandleTopGames(w http.ResponseWriter, r *http.Request) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:       "...",
		AppAccessToken: "...",
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.GetTopGames(&helix.TopGamesParams{
		First: 20,
	})
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(resp)
}

// HandleTopChannels responds with the top channels
func (a *application) HandleTopChannels(w http.ResponseWriter, r *http.Request) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:       "...",
		AppAccessToken: "...",
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.GetStreams(&helix.StreamsParams{
		First:    20,
		Language: []string{"en"},
	})
	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(resp)
}
