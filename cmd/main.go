package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gerbenjacobs/svc/handler"
	"github.com/gerbenjacobs/svc/services"
	"github.com/gerbenjacobs/svc/storage"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// handle shutdown signals
	shutdown := make(chan os.Signal, 3)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// set output logging (specifically for windows)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	log.SetLevel(log.DebugLevel)

	// load configuration
	var c Configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	// set up and check database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", c.DB.User, c.DB.Password, c.DB.Database))
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// create repositories and services
	auth := services.NewAuth([]byte(c.App.SecretToken))
	userSvc, err := services.NewUserSvc(storage.NewUserRepository(db), auth)
	if err != nil {
		log.Fatalf("failed to start user service: %v", err)
	}
	webhookSvc := services.NewWebhookService(storage.NewWebhookRepository(db))

	// set up the route handler and server
	mux := handler.New(handler.Dependencies{
		Auth:       auth,
		UserSvc:    userSvc,
		WebhookSvc: webhookSvc,
	})
	srv := &http.Server{
		Addr:         c.App.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	// start running the server
	go func() {
		log.Print("Server started")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen: %v", err)
		}
	}()

	// wait for shutdown signals
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Print("Server stopped successfully")
}

type Configuration struct {
	App struct {
		Address     string
		SecretToken string
	}
	DB struct {
		User     string
		Password string
		Database string
	}
}
