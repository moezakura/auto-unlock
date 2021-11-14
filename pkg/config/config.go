package config

type Config struct {
	UnlockToken    string   `config:"unlock_token,short=t" json:"unlock_token" yaml:"unlock_token"`
	SoundLevel     int      `config:"sound_level" json:"sound_level" yaml:"sound_level"`
	HostMacAddress []string `config:"host_mac_address" json:"host_mac_address" yaml:"host_mac_address"`
}
