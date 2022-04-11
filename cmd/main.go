package main

import (
	"context"
	"errors"
	"github.com/pmokeev/rate-limiter/tree/main/internal"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.ReadInConfig()
}

func main() {
	router := internal.NewRouter()
	server := internal.NewServer()

	go func() {
		if err := server.Run(viper.GetString("port"), router.InitRouter()); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	log.Println("API started")

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("API forced to shutdown:", err)
	}

	log.Println("API exiting")
}
