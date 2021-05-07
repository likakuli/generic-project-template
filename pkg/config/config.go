package config

type Config struct {
	Server *ServerConfig `mapstructure:"server"`
	Common *CommonConfig `mapstructure:"common"`
	DB     *DBConfig     `mapstructure:"db"`
}

type CommonConfig struct {
	DebugLevel int
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	API_QPS int `mapstructure:"api_qps"`
}

type DBConfig struct {
	ConnectionString string `mapstructure:"connectionString"`
	MaxOpenConn      int    `mapstructure:"maxOpenConn"`
	MaxIdleConn      int    `mapstructure:"maxIdleConn"`
}
