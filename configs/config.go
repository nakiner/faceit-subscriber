package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nakiner/faceit/pkg/store/nats"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
)

// ServiceName Used to define service prefix
const ServiceName = "faceit-subscriber"

// options slice to map all values into configuration
var options = []option{
	{"config", "string", "", "config file"},

	{"nats.host", "string", "127.0.0.1", "The nats host"},
	{"nats.port", "int", 4222, "The nats port"},
	{"nats.username", "string", "", "The nats user login"},
	{"nats.password", "string", "", "The nats user password"},
	{"nats.request_timeout_msec", "int", 500000, "The nats connection timeout in msec"},
	{"nats.retry_limit", "int", 5, "Reconnection limit to the nats"},
	{"nats.reconnect_time_wait_msec", "int", 500, "Reconnect time wait to the nats in msec"},

	{"logger.level", "string", "emerg", "Level of logging. A string that correspond to the following levels: emerg, alert, crit, err, warning, notice, info, debug"},
	{"logger.time_format", "string", "2006-01-02T15:04:05.999999999", "Date format in logs"},
}

type Config struct {
	Logger struct {
		Level      string
		TimeFormat string `mapstructure:"time_format"`
	}
	Nats nats.Config
}

type option struct {
	name        string
	typing      string
	value       interface{}
	description string
}

// NewConfig returns and prints struct with config parameters
func NewConfig() *Config {
	return &Config{}
}

// Read gets parameters from environment variables, flags or file.
func (c *Config) Read() error {
	viper.SetEnvPrefix(ServiceName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	for _, o := range options {
		switch o.typing {
		case "string":
			pflag.String(o.name, o.value.(string), o.description)
		case "int":
			pflag.Int(o.name, o.value.(int), o.description)
		case "bool":
			pflag.Bool(o.name, o.value.(bool), o.description)
		case "float64":
			pflag.Float64(o.name, o.value.(float64), o.description)
		default:
			viper.SetDefault(o.name, o.value)
		}
	}

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	viper.BindPFlags(pflag.CommandLine)
	pflag.Parse()

	if fileName := viper.GetString("config"); fileName != "" {
		viper.SetConfigFile(fileName)
		viper.SetConfigType("toml")

		if err := viper.ReadInConfig(); err != nil {
			return errors.Wrap(err, "failed to read from file")
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

// Print prints actual config on runtime start
func (c *Config) Print() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(b))
	return nil
}
