package config

type Config struct {
	UnlockToken string `config:"unlock_token,short=t" json:"unlock_token" yaml:"unlock_token"`
}
