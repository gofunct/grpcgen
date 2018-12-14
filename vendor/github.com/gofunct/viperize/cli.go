package viperize

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
)

func ViperizeCli() {
	viper.SetConfigName("viperize")         // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AutomaticEnv()
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache") // read in environment variables that match
	viper.SetDefault("root-usage", "example-usage") // read in environment variables that match
	viper.SetDefault("root-info", "replace with a short description of your cli application") // read in environment variables that match
	viper.SetDefault("root-summary", "replace with a long description of your cli application") // read in environment variables that match


	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := viper.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to --> viperize.yaml")
		if err := viper.WriteConfigAs("viperize.yaml"); err != nil {
			log.Fatal("failed to write config")
			os.Exit(1)
		}

	} else {
		log.Info("Using config file-->", viper.ConfigFileUsed())
		if err := viper.WriteConfig(); err != nil {
			log.Fatal("failed to write config file")
			os.Exit(1)
		}
	}
}
