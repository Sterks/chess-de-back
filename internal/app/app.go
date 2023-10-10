package app

import (
	"chess-backend/internal/config"
	"chess-backend/internal/delivery"
	"chess-backend/internal/repository"
	"chess-backend/internal/server"
	"chess-backend/internal/service"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

func Run(configPath string) error {

	cfg, err := config.InitConfig(configPath)
	if err != nil {
		log.Printf("Не могу получить данные из конфигурационного файла", err)
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoDB.URL))
	if err != nil {
		log.Println("Mongo is not running")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln("Does't work mongo server! ", err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalln("Mongo is not running")
		}
	}()

	repos := repository.NewRepositories(client)
	services := service.NewServices(service.Deps{
		Repos: repos,
	})
	handlers := delivery.NewHandler(services)

	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); err != nil {
			log.Printf("error occurred while running http server: %s\n", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("failed to stop server: %v", err)
	}
	return nil
}
