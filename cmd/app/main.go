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

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
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

	// set stackdriver formatter
	if c.Svc.Env != "dev" {
		log.SetLevel(log.InfoLevel)
		log.SetFormatter(stackdriver.NewFormatter(
			stackdriver.WithService(c.Svc.Name),
			stackdriver.WithVersion("v"+c.Svc.Version),
		))
	}

	// set up and check database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// create repositories and services
	auth := services.NewAuth([]byte(c.Svc.SecretToken))
	userSvc, err := services.NewUserSvc(storage.NewUserRepository(db), auth)
	if err != nil {
		log.Fatalf("failed to start user service: %v", err)
	}
	webhookSvc := services.NewWebhookService(storage.NewWebhookRepository(db))

	// set up the route handler and server
	app := handler.New(handler.Dependencies{
		Auth:       auth,
		UserSvc:    userSvc,
		WebhookSvc: webhookSvc,
	})
	srv := &http.Server{
		Addr:         c.Svc.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      app,
	}

	// start running the server
	go func() {
		log.Print("Server started on " + srv.Addr)
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
	Svc struct {
		Name        string
		Version     string
		Env         string
		Address     string
		SecretToken string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
}
