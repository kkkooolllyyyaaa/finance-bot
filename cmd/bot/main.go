package main

import (
	"log"

	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/clients/tg"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/command/commands"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/command"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/messages"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/model/reader"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/readers"
	repo "gitlab.ozon.dev/kolya_cypandin/project-base/internal/repository/expense"
	service "gitlab.ozon.dev/kolya_cypandin/project-base/internal/service/expense"
)

func main() {
	yamlReader := readers.New()
	readerModel := reader.New(yamlReader)
	secretManager, err := config.NewSecretManager(readerModel)
	if err != nil {
		log.Fatal("Can't create new SecretManager:", err)
	}

	expenseRepo := repo.NewInMemRepository()
	expenseMgmtService := service.NewMgmtService(expenseRepo)

	registry := command.NewRegistry()
	initCommands(registry, expenseMgmtService)

	token := secretManager.Token()
	configManager, err := config.NewConfigManager(readerModel)
	if err != nil {
		log.Fatal("Can't create new ConfigManager:", err)
	}
	tgClient, err := tg.New(token, registry, configManager)

	if err != nil {
		log.Fatal("Can't create new tgClient:", err)
	}

	msgModel := messages.New(tgClient)
	tgClient.ListenUpdates(msgModel)
}

func initCommands(registry *command.Registry, expenseMgmtService *service.MgmtService) {
	registry.AddCommand(
		"/add",
		commands.NewAddExpenseCommand(expenseMgmtService),
	)
	registry.AddCommand(
		"/report",
		commands.NewReportCommand(expenseMgmtService),
	)
	registry.AddCommand(
		"/help",
		commands.NewHelpCommand(
			registry.Commands(),
		),
	)
}
