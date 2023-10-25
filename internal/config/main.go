package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type SqliteConf struct {
	User     string
	Pass     string
	Database string
}

type MariadbConf struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type PostgresConf struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type DbConf struct {
	Driver       string `mapstructure:"engine"`
	SqliteConf   `mapstructure:"sqlite"`
	MariadbConf  `mapstructure:"mariadb"`
	PostgresConf `mapstructure:"postgres"`
}

type ServerConf struct {
	Env  string
	Port string
}

type Config struct {
	*DbConf     `mapstructure:"database"`
	*ServerConf `mapstructure:"server"`
	vpr         *viper.Viper
}

func (c *Config) LoadEnvVariables() {
	c.vpr = viper.GetViper()
	c.vpr.SetConfigType("yaml")
	c.vpr.SetConfigName("config")
	for _, envPath := range []string{"./", "../", "../../"} {
		c.vpr.AddConfigPath(envPath)
	}

	if err := c.vpr.ReadInConfig(); err != nil {
		log.Fatalf("fail to read config file, err: %+v\n", err)
	}

	/* Set default */
	c.vpr.SetDefault("server", map[string]string{
		"env":  "local",
		"port": "7000",
	})

	// server := c.vpr.Get("server")
	// fmt.Printf("server: %+v\n", server)
	// database := c.vpr.Get("database")
	// fmt.Printf("database: %+v\n", database)

	// load server settings
	if err := c.vpr.Unmarshal(c); err != nil {
		log.Printf("failed loading conf, err: %+v\n", err.Error())
	}
	// log.Printf("conf: %+v\n", *c.ServerConf)
	// log.Printf("conf: %+v\n", *c.DbConf)
	log.Printf("loaded config.yaml successfully")
}

func (c *Config) WatchConfig() {
	c.vpr.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})

	c.vpr.WatchConfig()
}
