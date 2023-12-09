package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type SqliteConf struct {
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type MariadbConf struct {
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type PostgresConf struct {
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type MongodbConf struct {
	Host     *string
	Port     *string
	User     *string
	Pass     *string
	Database *string
}

type DbConf struct {
	Driver       string `mapstructure:"engine"`
	SqliteConf   `mapstructure:"sqlite"`
	MariadbConf  `mapstructure:"mariadb"`
	PostgresConf `mapstructure:"postgres"`
	MongodbConf `mapstructure:"mongodb"`
}

type ServerConf struct {
	Env  string
	Host string
	Port string
}

type Jwt struct {
	Secret string
}

type Config struct {
	*DbConf     `mapstructure:"database"`
	*ServerConf `mapstructure:"server"`
	*Jwt        `mapstructure:"jwt"`
	Vpr         *viper.Viper
}

func (c *Config) LoadEnvVariables() {
	c.Vpr = viper.GetViper()
	c.Vpr.SetConfigType("yaml")
	c.Vpr.SetConfigName("config")
	for _, envPath := range []string{"./", "../", "../../"} {
		c.Vpr.AddConfigPath(envPath)
	}

	if err := c.Vpr.ReadInConfig(); err != nil {
		log.Fatalf("fail to read config file, err: %+v\n", err)
	}

	/* Set default */
	c.Vpr.SetDefault("server", map[string]string{
		"env":  "local",
		"port": "7000",
	})

	// server := c.Vpr.Get("server")
	// log.Printf("server: %+v\n", server)
	// database := c.Vpr.Get("database")
	// log.Printf("database: %+v\n", database)

	// load server settings
	if err := c.Vpr.Unmarshal(c); err != nil {
		log.Printf("failed loading conf, err: %+v\n", err.Error())
	}
	// log.Printf("conf: %+v\n", *c.ServerConf)
	// log.Printf("conf: %+v\n", *c.DbConf)
	log.Printf("loaded config.yaml successfully")
}

func (c *Config) WatchConfig() {
	c.Vpr.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		c.LoadEnvVariables()
	})

	c.Vpr.WatchConfig()
}

var Cfg = Config{}
