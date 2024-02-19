package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
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
	MongodbConf  `mapstructure:"mongodb"`
}

type Logging struct {
	Level int
	Type  []string
	Zap   struct {
		Output   []string
		Filename string
	}
	DebugSymbol *string
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
	*Logging    `mapstructure:"logging"`
	Vpr         *viper.Viper
}

func (c *Config) LoadEnvVariables() {
	c.Vpr.SetConfigType("yaml")

	// determine the /.dockerenv file for checking running inside docker or not for using the corresponding config
	// ref: https://stackoverflow.com/a/12518877
	if _, err := os.Stat("/.dockerenv"); err == nil { // running in docker
		// path/to/whatever exists
		// log.Printf("Running inside docker\n")
		c.Vpr.SetConfigName("config-docker")
	} else if errors.Is(err, os.ErrNotExist) { // running in localhost w/o docker
		// path/to/whatever does *not* exist
		// log.Printf("Running in localhost\n")
		c.Vpr.SetConfigName("config-localhost")
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Printf("env check for config err: %+v\n", err)
	}

	for _, envPath := range []string{"./configs"} {
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

var Cfg = &Config{
	Vpr: viper.GetViper(),
}
