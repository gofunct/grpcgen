package runtime

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func NewGrpcServerViperizer(vi *viper.Viper, service string) error {
	var cfg = "grpcserver.yaml"
	vi.SetConfigName("grpcserver")       // name of config file (without extension)
	vi.AddConfigPath(os.Getenv("$HOME")) // name of config file (without extension)
	vi.AddConfigPath(".")                // optionally look for config in the working directory
	vi.AutomaticEnv()
	vi.SetEnvPrefix(service + "_SERVER")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.SetDefault("proxy.service", service)
	vi.SetDefault("server.tls", false)
	vi.SetDefault("server.tracing", true)
	vi.SetDefault("server.metrics_endpoint", true)
	vi.SetDefault("server.live_endpoint", true)
	vi.SetDefault("server.ready_endpoint", true)
	vi.SetDefault("server.pprof_endpoint", true)
	vi.SetDefault("server.db_host", "localhost")
	vi.SetDefault("server.db_port", ":5432")
	vi.SetDefault("server.db_name", "postgresdb")
	vi.SetDefault("server.db_user", "admin")
	vi.SetDefault("server.port", ":8443")
	vi.SetDefault("server.routine_threshold", 300)
	vi.SetDefault("server.jaeger_metrics", true)

	// If a config file is found, read it in."A generator for gRPC based Applications"
	if err := vi.ReadInConfig(); err != nil {
		log.Info("failed to read config file, writing defaults to -->" + cfg)
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


