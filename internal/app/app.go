package app

import (
	"context"
	"sync"
	"time"

	currencyApi "gitlab.ozon.dev/kolya_cypandin/project-base/internal/client/currency"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/config/inmem_config"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/command/commands"
	commandService "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/command/service"
	currencyService "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/currency/service"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/entity"
	expenseRepository "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/expense/repository"
	expenseService "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/expense/service"
	messagingService "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/messaging/service"
	userRepository "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user/repository"
	userService "gitlab.ozon.dev/kolya_cypandin/project-base/internal/domain/user/service"
	"gitlab.ozon.dev/kolya_cypandin/project-base/internal/util/file/yaml"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, erg *errgroup.Group) error {
	fileReader := yaml.New()

	configFilePath := "configs/application.yaml"
	configManager, err := inmem_config.NewConfigManager(fileReader, configFilePath)
	if err != nil {
		return err
	}

	secretsFilePath := "configs/secrets.yaml"
	secretManager, err := inmem_config.NewSecretManager(fileReader, secretsFilePath)
	if err != nil {
		return err
	}

	currencyApiImpl := new(currencyApi.Api)
	currencyServiceImpl := currencyService.New(currencyApiImpl)
	erg.Go(func() error {
		updateTickSeconds := configManager.CurrencyRatesUpdateTick()
		ticker := time.NewTicker(time.Duration(updateTickSeconds) * time.Second)

		return currencyServiceImpl.ListenForUpdates(ctx, ticker)
	})

	expenseRepositoryImpl := expenseRepository.NewInMemRepository()
	expenseMgmtService := expenseService.NewMgmtService(expenseRepositoryImpl)

	userRepositoryImpl := userRepository.NewInMemRepository()
	userMgmtService := userService.NewMgmtService(userRepositoryImpl)

	commandRegistry := commandService.NewRegistry()
	initCommands(commandRegistry, expenseMgmtService, userMgmtService, currencyServiceImpl)

	token := secretManager.Token()
	tgMessagingService, err := messagingService.NewTgMessaging(token, userMgmtService, configManager)
	if err != nil {
		return err
	}

	bufferSize := configManager.MessagingBufferSize()
	gotMessagesChan := make(chan entity.Message, bufferSize)
	toSendMessagesChan := make(chan entity.Message, bufferSize)

	erg.Go(func() error {
		err = tgMessagingService.ListenForMessages(ctx, gotMessagesChan)
		close(gotMessagesChan)
		return err
	})

	workersCount := configManager.MessagingWorkersCount()
	var once sync.Once
	onceBody := func() { close(toSendMessagesChan) }
	for i := 0; i < workersCount; i++ {
		erg.Go(func() error {
			err = commandRegistry.ExecuteCommands(ctx, gotMessagesChan, toSendMessagesChan)
			once.Do(onceBody)
			return err
		})
	}

	erg.Go(func() error {
		return tgMessagingService.SendMessages(ctx, erg, toSendMessagesChan)
	})

	return nil
}

func initCommands(
	registry *commandService.CommandService,
	expenseMgmtService *expenseService.MgmtService,
	userMgmtService *userService.MgmtService,
	currencyRatesStorage *currencyService.CurrencyRatesStorage,
) {
	registry.AddCommand(
		"/add",
		commands.NewAddExpenseCommand(expenseMgmtService, userMgmtService, currencyRatesStorage),
	)
	registry.AddCommand(
		"/report",
		commands.NewReportCommand(expenseMgmtService, userMgmtService, currencyRatesStorage),
	)
	registry.AddCommand(
		"/help",
		commands.NewHelpCommand(
			registry.Commands(),
		),
	)
	registry.AddCommand(
		"/start",
		commands.NewStartCommand(),
	)
	registry.AddCommand(
		"/currency",
		commands.NewSetCurrencyCommand(userMgmtService),
	)
}
