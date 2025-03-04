package config

import (
	"errors"
	"fmt"
	"golang-api-starter/internal/helper/utils"
	"log"
	"os"

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

type RabbitMqConf struct {
	Host   *string
	Port   *string
	User   *string
	Pass   *string
	Queues struct {
		LogQueue   *string
		EmailQueue *string
		TestQueue  *string
	}
}

type RedisConf struct {
	Host string
	Port string
	User *string
	Pass *string
}
type MemcachedConf struct {
	Host string
	Port string
	User *string
	Pass *string
}
type CacheConf struct {
	Enabled       bool
	Driver        string `mapstructure:"engine"`
	RedisConf     `mapstructure:"redis"`
	MemcachedConf `mapstructure:"memcached"`
}

type ApsaraConf struct {
	AccessKey    *string
	AccessSecret *string
	PushKey      *string
	PullKey      *string
}

type TranscodingApi struct {
	Host   *string
	Port   *string
	Secure bool
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
	Env            string
	Host           string
	Port           string
	TrustedProxies []string
}

type Jwt struct {
	Secret string
}

type OAuthGoogle struct {
	Key         string
	Secret      string
	CallbackUrl string
}
type OAuthGithub struct {
	Key         string
	Secret      string
	CallbackUrl string
}
type OAuth struct {
	*OAuthGoogle `mapstructure:"google"`
	*OAuthGithub `mapstructure:"github"`
}

type Smtp struct {
	Host string
	Port int
	Ssl  bool
	User string
	Pass string
}
type Notification struct {
	Smtp *Smtp
}

type Config struct {
	*DbConf       `mapstructure:"database"`
	*ServerConf   `mapstructure:"server"`
	*Jwt          `mapstructure:"jwt"`
	*Logging      `mapstructure:"logging"`
	*OAuth        `mapstructure:"oauth"`
	*Notification `mapstructure:"notification"`
	*RabbitMqConf `mapstructure:"rbmq"`
	*CacheConf    `mapstructure:"cache"`
	Vpr           *viper.Viper
}

// LoadEnvVariables loads the config yaml file from ./configs/
func (c *Config) LoadEnvVariables() {
	c.Vpr.SetConfigType("yaml")

	// determine the /.dockerenv file for checking running inside docker or not for using the corresponding config
	// ref: https://stackoverflow.com/a/12518877
	if _, err := os.Stat("/.dockerenv"); err == nil { // running in docker
		log.Printf("Running inside docker\n")
		c.Vpr.SetConfigName("docker")
		c.setDockerDefault()
	} else if len(os.Getenv("KUBERNETES_SERVICE_HOST")) > 0 { // running in k8s ref: https://stackoverflow.com/a/54130803
		log.Printf("Running in k8s\n")
		c.Vpr.SetConfigName("k3s")
		c.setK3sDefault()
	} else if errors.Is(err, os.ErrNotExist) { // running in localhost w/o docker
		log.Printf("Running in localhost\n")
		c.Vpr.SetConfigName("localhost")
		c.setLocalDefault()
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		log.Fatalf("env check for config err: %+v\n", err)
	}

	basepath := utils.RootDir(2)
	configsDir := fmt.Sprintf("%s/configs", basepath)
	// log.Printf("configsDir: %+v\n\n", configsDir)
	for _, envPath := range []string{configsDir} {
		c.Vpr.AddConfigPath(envPath)
	}

	if err := c.Vpr.ReadInConfig(); err != nil {
		log.Fatalf("fail to read config file, err: %+v\n", err)
	}

	// server := c.Vpr.Get("server")
	// log.Printf("server: %+v\n", server)
	// database := c.Vpr.Get("database")
	// log.Printf("database: %+v\n", database)

	// load server settings
	if err := c.Vpr.Unmarshal(c); err != nil {
		log.Fatalf("failed loading conf, err: %+v\n", err.Error())
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

// GetServerUrl returns server url by config
func (c *Config) GetServerUrl() string {
	url := fmt.Sprintf("http://%s", c.ServerConf.Host)

	if len(c.ServerConf.Port) > 0 && c.Env == "local" {
		url = fmt.Sprintf("%s:%s", url, c.ServerConf.Port)
	}

	return url
}

var Cfg = &Config{
	Vpr: viper.GetViper(),
}
