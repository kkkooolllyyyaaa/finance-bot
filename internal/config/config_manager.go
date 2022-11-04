package config

type ConfigManager interface {
	UpdateTimeout() int
	MessagingBufferSize() int
	MessagingRetryCount() int
	MessagingWorkersCount() int
	CurrencyRatesUpdateTick() int
}
