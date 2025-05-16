package config

type Config struct {
	InfluxDB  *string
	Blacklist *[]string
}

var config = &Config{
	InfluxDB:  nil,
	Blacklist: nil,
}

func GetConfig() *Config {
	return config
}
