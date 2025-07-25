package cfgmanager

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type ConfigManager struct {
	variables []string
}

func NewConfigManager(variables []string) (*ConfigManager, error) {
	if err := setSystemVars(variables); err != nil {
		return nil, err
	}
	return &ConfigManager{
		variables: variables,
	}, nil
}

func (c *ConfigManager) ClearVars() {
	for _, variable := range c.variables {
		os.Unsetenv(variable)
	}
}

func setSystemVars(variables []string) error {
	// Set the file name of the configurations file
	viper.SetConfigName("config.yaml")

	// Set the path to look for the configurations file
	viper.AddConfigPath("../configs/")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file", "error", err.Error())
		return err
	}

	for _, varName := range variables {
		value, ok := viper.Get(varName).(string)
		if !ok {
			slog.Error("Invalid type assertion", "variable name", varName)
		}
		if err := os.Setenv(varName, value); err != nil {
			return err
		}
	}

	return nil
}
