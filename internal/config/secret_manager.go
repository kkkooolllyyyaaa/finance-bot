package config

type SecretManager interface {
	Token() string
}
