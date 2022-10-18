package main

import (
	"log"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/clients/tg"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config/reader"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/messages"
)

func main() {
	configReader := reader.New()
	secretManager, err := config.New(configReader)
	if err != nil {
		log.Fatal("Can't create new Secret:", err)
	}
	token := secretManager.Token()

	tgClient, err := tg.New(token)
	if err != nil {
		log.Fatal("Can't create new tgClient:", err)
	}

	msgModel := messages.New(tgClient)
	tgClient.ListenUpdates(msgModel)
}
