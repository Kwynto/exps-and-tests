package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/Kwynto/read-adviser-bot/clients/telegram"
	event_consumer "github.com/Kwynto/read-adviser-bot/consumer/event-consumer"
	"github.com/Kwynto/read-adviser-bot/events/telegram"
	"github.com/Kwynto/read-adviser-bot/storage/sqlite"
)

const (
	tgBotHost          = "api.telegram.org"
	storagePath        = "files_storage"
	storageStoragePath = "data/sqlite/storage.db"
	batchSize          = 100
)

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}

func main() {
	// s := files.New(storagePath)
	s, err := sqlite.New(storageStoragePath)
	if err != nil {
		log.Fatalf("can't connect to storage: %s", err)
	}

	err = s.Init(context.Background())
	if err != nil {
		log.Fatalf("can't init storage: %s", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}
