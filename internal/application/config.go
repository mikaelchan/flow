package application

type PostgreConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type ServiceConfig struct {
	Postgre PostgreConfig `toml:"postgre"`
	Redis   RedisConfig   `toml:"redis"`
}
