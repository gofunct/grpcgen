package proxy

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewProxyViperizer(vi *viper.Viper, service string) error {
	var cfg = "proxy.yaml"
	vi.SetConfigName("proxy") // name of config file (without extension)
	vi.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	vi.AddConfigPath(".")                // optionally look for config in the working directory
	vi.AutomaticEnv()
	vi.SetEnvPrefix(service+"_PROXY")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.SetDefault("proxy.service", service)
	vi.SetDefault("proxy.prefix", "")
	vi.SetDefault("proxy.tracing", true)
	vi.SetDefault("proxy.tls", false)
	vi.SetDefault("proxy.metrics_endpoint", true)
	vi.SetDefault("proxy.live_endpoint", true)
	vi.SetDefault("proxy.ready_endpoint", true)
	vi.SetDefault("proxy.pprof_endpoint", true)
	vi.SetDefault("proxy.db_host", "localhost")
	vi.SetDefault("proxy.db_port", ":5432")
	vi.SetDefault("proxy.db_name", "postgresdb")
	vi.SetDefault("proxy.db_user", "admin")
	vi.SetDefault("proxy.routine_threshold", 300)
	vi.SetDefault("proxy.jaeger_metrics", true)
	vi.SetDefault("proxy.port", ":8080")
	vi.SetDefault("proxy.backend", ":8443")
	vi.SetDefault("proxy.log_level", "")
	vi.SetDefault("proxy.swagger_file", "")
	vi.SetDefault("proxy.allow_origin", "")
	vi.SetDefault("proxy.allow_creds", "")
	vi.SetDefault("proxy.allow_methods", "")
	vi.SetDefault("proxy.allow_headers", "")

	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := vi.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->"+cfg)
		if err := vi.WriteConfigAs(cfg); err != nil {
			return err
		}

	} else {
		log.Info("Using config file-->", vi.ConfigFileUsed())
		if err := vi.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}


