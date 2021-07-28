package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config the asd.json file should be set at the root level
type Config struct {
	SERVER struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	DB struct {
		Driver      string `json:"driver" required:"true"`
		URL         string `json:"url" required:"true"`
		Automigrate bool
	}
}

var cfgFile string

func LoadTestConfig(path string) (*Config, error) {
	cmd := &cobra.Command{}
	cmd.Flags().StringVar(&cfgFile, "config", path, "Config file")

	return LoadConfig(cmd)
}

// LoadConfig should load and unmarshal the config file
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("GOKIT")

	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	config := new(Config)

	if err := viper.Unmarshal(&config); err != nil {
		zap.L().Sugar().Fatalf("Unable to decode into struct")
	}

	fmt.Println(config)

	return config, nil
}
